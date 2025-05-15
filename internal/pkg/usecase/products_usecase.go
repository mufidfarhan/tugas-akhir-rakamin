package usecase

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/pkg/entity"
	"backend-evermos/internal/pkg/model"
	"backend-evermos/internal/pkg/repository"
	"backend-evermos/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductsUseCase interface {
	CreateProduct(ctx context.Context, userID string, data model.ProductReqCreate, files []*multipart.FileHeader) (res uint, err *helper.ErrorStruct)
	GetAllProducts(ctx context.Context, params model.ProductsFilter) (res model.FilteredData, err *helper.ErrorStruct)
	GetProductByID(ctx context.Context, productID string) (res model.ProductResp, err *helper.ErrorStruct)
	UpdateProductByID(ctx context.Context, userID string, productID string, data model.ProductReqUpdate, files []*multipart.FileHeader) (res string, err *helper.ErrorStruct)
	DeleteProductByID(ctx context.Context, userID string, productID string) (res string, err *helper.ErrorStruct)
}

type ProductsUseCaseImpl struct {
	productsRepository      repository.ProductsRepository
	shopsRepository         repository.ShopsRepository
	productImagesRepository repository.ProductImagesRepository
	categoriesRepository    repository.CategoriesRepository
}

func NewProductsUseCase(
	productsRepository repository.ProductsRepository,
	shopsRepository repository.ShopsRepository,
	productImagesRepository repository.ProductImagesRepository,
	categoriesRepository repository.CategoriesRepository,
) ProductsUseCase {
	return &ProductsUseCaseImpl{
		productsRepository:      productsRepository,
		shopsRepository:         shopsRepository,
		productImagesRepository: productImagesRepository,
		categoriesRepository:    categoriesRepository,
	}
}

func (alc *ProductsUseCaseImpl) CreateProduct(ctx context.Context, userID string, data model.ProductReqCreate, files []*multipart.FileHeader) (res uint, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	var shopID uint
	resRepo, errRepo := alc.shopsRepository.GetShopByUserID(ctx, userID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("toko tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetShopByUserID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	categoryID := utils.NilIfZeroUint(data.CategoryID)
	if categoryID != nil {
		categoryID := fmt.Sprintf("%d", *data.CategoryID)
		if errRepo := alc.categoriesRepository.VerifyCategoryAvailability(ctx, categoryID); errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return res, &helper.ErrorStruct{
					Code: fiber.StatusNotFound,
					Err:  errors.New("kategori tidak ditemukan"),
				}
			}

			helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyCategoryAvailability: %s", errRepo.Error()), errRepo)
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errRepo,
			}
		}
	}

	uploadDir := "files/products"
	var photoURLs []string
	for _, fileHeader := range files {
		photoURL, err := utils.SaveFileToDisk(fileHeader, uploadDir)
		if err != nil {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  err,
			}
		}
		photoURLs = append(photoURLs, photoURL)
	}

	var productID uint
	errTransaction := alc.productsRepository.WithinTransaction(ctx, func(txCtx context.Context) (err error) {
		shopID = resRepo.ID
		productID, err = alc.productsRepository.CreateProduct(txCtx, entity.Product{
			ProductName:   data.ProductName,
			Slug:          data.Slug,
			ResellerPrice: data.ResellerPrice,
			ConsumerPrice: data.ConsumerPrice,
			Stock:         data.Stock,
			Description:   data.Description,
			ShopID:        shopID,
			CategoryID:    categoryID,
		})
		if err != nil {
			return err
		}

		fmt.Println(productID)

		for _, photoURL := range photoURLs {
			_, err = alc.productImagesRepository.CreateProductImage(txCtx, entity.ProductImage{
				ProductID: productID,
				PhotoURL:  photoURL,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		for _, photoURL := range photoURLs {
			_ = os.Remove(photoURL)
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at WithinTransaction: %s", errTransaction.Error()), errTransaction)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errors.New("gagal menambahkan produk"),
		}
	}

	return productID, nil
}

func (alc *ProductsUseCaseImpl) GetAllProducts(ctx context.Context, params model.ProductsFilter) (res model.FilteredData, err *helper.ErrorStruct) {
	var products []model.ProductResp

	limit, offset := func(limit, page int) (int, int) {
		if limit < 1 {
			limit = 10
		}

		var offset int
		if page < 1 {
			offset = 0
		} else {
			offset = (page - 1) * limit
		}
		return limit, offset
	}(params.Limit, params.Page)

	resRepo, errRepo := alc.productsRepository.GetAllProducts(ctx, entity.FilterProducts{
		Limit:       limit,
		Offset:      offset,
		ProductName: params.ProductName,
		CategoryID:  params.CategoryID,
		ShopID:      params.ShopID,
		MaxPrice:    params.MaxPrice,
		MinPrice:    params.MinPrice,
	})

	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("produk tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAllProducts: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	for _, v := range resRepo {
		var imageResponses []model.ProductImageResp
		for _, img := range v.Images {
			imageResponses = append(imageResponses, model.ProductImageResp{
				ID:        img.ID,
				ProductID: img.ProductID,
				ImageURL:  img.PhotoURL,
			})
		}

		products = append(products, model.ProductResp{
			ID:            v.ID,
			ProductName:   v.ProductName,
			Slug:          v.Slug,
			ResellerPrice: v.ResellerPrice,
			ConsumerPrice: v.ConsumerPrice,
			Stock:         v.Stock,
			Description:   v.Description,
			Shop: model.ShopResp{
				ID:       v.Shop.ID,
				ShopName: v.Shop.ShopName,
				PhotoURL: v.Shop.PhotoURL,
			},
			Category: model.CategoryResp{
				ID:           v.Category.ID,
				CategoryName: v.Category.CategoryName,
			},
			Images: imageResponses,
		})
	}

	res = model.FilteredData{
		Data:  products,
		Page:  params.Page,
		Limit: params.Limit,
	}

	return res, err
}

func (alc *ProductsUseCaseImpl) GetProductByID(ctx context.Context, productID string) (res model.ProductResp, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.productsRepository.GetProductByID(ctx, productID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("produk tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetProductByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	var images []model.ProductImageResp
	for _, img := range resRepo.Images {
		images = append(images, model.ProductImageResp{
			ID:        img.ID,
			ProductID: img.ProductID,
			ImageURL:  img.PhotoURL,
		})
	}

	res = model.ProductResp{
		ID:            resRepo.ID,
		ProductName:   resRepo.ProductName,
		Slug:          resRepo.Slug,
		ResellerPrice: resRepo.ResellerPrice,
		ConsumerPrice: resRepo.ConsumerPrice,
		Stock:         resRepo.Stock,
		Description:   resRepo.Description,
		Shop: model.ShopResp{
			ID:       resRepo.Shop.ID,
			ShopName: resRepo.Shop.ShopName,
			PhotoURL: resRepo.Shop.PhotoURL,
		},
		Category: model.CategoryResp{
			ID:           resRepo.Category.ID,
			CategoryName: resRepo.Category.CategoryName,
		},
		Images: images,
	}

	return res, err
}

func (alc *ProductsUseCaseImpl) UpdateProductByID(ctx context.Context, userID string, productID string, data model.ProductReqUpdate, files []*multipart.FileHeader) (res string, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		log.Println(errValidate)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	resProductRepo, errRepo := alc.productsRepository.GetProductByID(ctx, productID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("produk tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetProductByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	resShopRepo, errRepo := alc.shopsRepository.GetShopByUserID(ctx, userID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("toko tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetShopByUserID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	shopID := fmt.Sprintf("%d", resShopRepo.ID)
	if errRepo := alc.productsRepository.VerifyProductOwner(ctx, productID, shopID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusForbidden,
				Err:  errors.New("anda tidak berhak mengakses resource ini"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyShopOwner: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	categoryID := utils.NilIfZeroUint(data.CategoryID)
	if categoryID != nil {
		categoryID := fmt.Sprintf("%d", *data.CategoryID)
		if errRepo := alc.categoriesRepository.VerifyCategoryAvailability(ctx, categoryID); errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return res, &helper.ErrorStruct{
					Code: fiber.StatusNotFound,
					Err:  errors.New("kategori tidak ditemukan"),
				}
			}

			helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyCategoryAvailability: %s", errRepo.Error()), errRepo)
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errRepo,
			}
		}
	}

	uploadDir := "files/products"
	var photoURLs []string
	for _, fileHeader := range files {
		photoURL, err := utils.SaveFileToDisk(fileHeader, uploadDir)
		if err != nil {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  err,
			}
		}
		photoURLs = append(photoURLs, photoURL)
	}

	errTransaction := alc.productsRepository.WithinTransaction(ctx, func(txCtx context.Context) (err error) {
		err = alc.productsRepository.UpdateProductByID(txCtx, productID, entity.Product{
			ProductName:   data.ProductName,
			Slug:          data.Slug,
			CategoryID:    categoryID,
			ResellerPrice: data.ResellerPrice,
			ConsumerPrice: data.ConsumerPrice,
			Stock:         data.Stock,
			Description:   data.Description,
		})
		if err != nil {
			return err
		}

		for i, photoURL := range photoURLs {
			if i >= len(photoURLs) {
				break
			}

			imageID := fmt.Sprintf("%d", resProductRepo.Images[i].ID)
			err = alc.productImagesRepository.UpdateProductImage(txCtx, imageID, entity.ProductImage{
				PhotoURL: photoURL,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		for _, photoURL := range photoURLs {
			_ = os.Remove(photoURL)
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at WithinTransaction: %s", errTransaction.Error()), errTransaction)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errors.New("gagal mengubah produk"),
		}
	}

	return "updated", nil
}

func (alc *ProductsUseCaseImpl) DeleteProductByID(ctx context.Context, userID string, productID string) (res string, err *helper.ErrorStruct) {
	resRepo, errRepo := alc.shopsRepository.GetShopByUserID(ctx, userID)
	if errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("toko tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetShopByUserID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	shopID := fmt.Sprintf("%d", resRepo.ID)
	if errRepo := alc.productsRepository.VerifyProductOwner(ctx, productID, shopID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusForbidden,
				Err:  errors.New("anda tidak berhak mengakses resource ini"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyShopOwner: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	if errRepo := alc.productsRepository.DeleteProductByID(ctx, productID); errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at DeleteProductByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	return "deleted", nil
}
