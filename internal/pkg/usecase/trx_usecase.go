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
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TrxUseCase interface {
	CreateTrx(ctx context.Context, userID string, data model.TrxReqCreate) (res uint, err *helper.ErrorStruct)
	GetTrxByID(ctx context.Context, userID string, trxID string) (res model.TrxResp, err *helper.ErrorStruct)
	GetAllTrx(ctx context.Context, userID string, params model.TrxFilter) (res model.FilteredData, err *helper.ErrorStruct)
}

type TrxUseCaseImpl struct {
	trxRepository           repository.TrxRepository
	trxDetailsRepository    repository.TrxDetailsRepository
	productLogsRepository   repository.ProductLogsRepository
	productsRepository      repository.ProductsRepository
	addressesRepository     repository.AddressesRepository
	productImagesRepository repository.ProductImagesRepository
}

func NewTrxUseCase(
	trxRepository repository.TrxRepository,
	trxDetailsRepository repository.TrxDetailsRepository,
	productLogsRepository repository.ProductLogsRepository,
	productsRepository repository.ProductsRepository,
	addressesRepository repository.AddressesRepository,
	productImagesRepository repository.ProductImagesRepository,
) TrxUseCase {
	return &TrxUseCaseImpl{
		trxRepository:           trxRepository,
		trxDetailsRepository:    trxDetailsRepository,
		productLogsRepository:   productLogsRepository,
		productsRepository:      productsRepository,
		addressesRepository:     addressesRepository,
		productImagesRepository: productImagesRepository,
	}
}

func (alc *TrxUseCaseImpl) CreateTrx(ctx context.Context, userID string, data model.TrxReqCreate) (res uint, err *helper.ErrorStruct) {
	if errValidate := helper.Validate.Struct(data); errValidate != nil {
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errValidate,
		}
	}

	addressID := fmt.Sprintf("%d", data.AddressID)
	if errRepo := alc.addressesRepository.VerifyAddressOwner(ctx, addressID, userID); errRepo != nil {
		if errors.Is(errRepo, gorm.ErrRecordNotFound) {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusNotFound,
				Err:  errors.New("alamat tidak ditemukan"),
			}
		}

		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at VerifyAddressOwner: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	var productTrx []entity.ProductTrx
	var grandTotal int

	for _, trxDetail := range data.TrxDetails {
		productID := fmt.Sprintf("%d", trxDetail.ProductID)
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

		if resRepo.Stock < trxDetail.Quantity {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("stok tidak tersedia"),
			}
		}

		price, err := strconv.Atoi(resRepo.ConsumerPrice)
		if err != nil {
			return res, &helper.ErrorStruct{
				Code: fiber.StatusInternalServerError,
				Err:  errors.New("gagal mengonversi harga ke integer"),
			}
		}

		productTotal := price * trxDetail.Quantity
		grandTotal += productTotal

		newStock := resRepo.Stock - trxDetail.Quantity

		productTrx = append(productTrx, entity.ProductTrx{
			Quantity:      trxDetail.Quantity,
			TotalPrice:    productTotal,
			ProductID:     resRepo.ID,
			ProductName:   resRepo.ProductName,
			Slug:          resRepo.Slug,
			ResellerPrice: resRepo.ResellerPrice,
			ConsumerPrice: resRepo.ConsumerPrice,
			Stock:         newStock,
			Description:   resRepo.Description,
			ShopID:        resRepo.ShopID,
			CategoryID:    resRepo.CategoryID,
		})
	}

	userIDNum, _ := utils.ConvertStringToUint(userID)
	invoiceCode := utils.GenerateInvoiceCode()

	var trxID uint
	errTransaction := alc.trxRepository.WithinTransaction(ctx, func(txCtx context.Context) (err error) {
		trxID, err = alc.trxRepository.CreateTrx(txCtx, entity.Trx{
			UserID:        userIDNum,
			AddressID:     data.AddressID,
			TotalPrice:    grandTotal,
			InvoiceCode:   invoiceCode,
			PaymentMethod: data.PaymentMethod,
		})
		if err != nil {
			return err
		}

		for _, data := range productTrx {
			productLogID, err := alc.productLogsRepository.CreateProductLogs(txCtx, entity.ProductLog{
				ProductID:     data.ProductID,
				ProductName:   data.ProductName,
				Slug:          data.Slug,
				ResellerPrice: data.ResellerPrice,
				ConsumerPrice: data.ConsumerPrice,
				Description:   data.Description,
				ShopID:        data.ShopID,
				CategoryID:    data.CategoryID,
			})
			if err != nil {
				return err
			}

			_, err = alc.trxDetailsRepository.CreateTrxDetails(txCtx, entity.TrxDetail{
				TrxID:        trxID,
				ProductLogID: productLogID,
				ShopID:       data.ShopID,
				Quantity:     data.Quantity,
				TotalPrice:   data.TotalPrice,
			})
			if err != nil {
				return err
			}

			productID := fmt.Sprintf("%d", data.ProductID)
			err = alc.productsRepository.UpdateProductByID(txCtx, productID, entity.Product{
				Stock: data.Stock,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at WithinTransaction: %s", errTransaction.Error()), errTransaction)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusInternalServerError,
			Err:  errors.New("gagal menambah transaksi"),
		}
	}

	return trxID, nil
}

func (alc *TrxUseCaseImpl) GetTrxByID(ctx context.Context, userID string, trxID string) (res model.TrxResp, err *helper.ErrorStruct) {
	trxResRepo, errRepo := alc.trxRepository.GetTrxByID(ctx, userID, trxID)
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

	addressID := fmt.Sprintf("%d", trxResRepo.AddressID)
	addressResRepo, errRepo := alc.addressesRepository.GetAddressByID(ctx, addressID)
	if !errors.Is(errRepo, gorm.ErrRecordNotFound) && errRepo != nil {
		helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAddressByID: %s", errRepo.Error()), errRepo)
		return res, &helper.ErrorStruct{
			Code: fiber.StatusBadRequest,
			Err:  errRepo,
		}
	}

	var trxDetails []model.TrxDetailResp
	for _, td := range trxResRepo.TrxDetails {

		productID := td.ProductLog.ProductID
		imgResRepo, errRepo := alc.productImagesRepository.GetImagesByProductID(ctx, productID)
		if errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return res, &helper.ErrorStruct{
					Code: fiber.StatusNotFound,
					Err:  errors.New("foto tidak ditemukan"),
				}
			}

			helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetImagesByProductID: %s", errRepo.Error()), errRepo)
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errRepo,
			}
		}

		var images []model.ProductImageResp
		for _, img := range imgResRepo {
			images = append(images, model.ProductImageResp{
				ID:        img.ID,
				ProductID: img.ProductID,
				ImageURL:  img.PhotoURL,
			})
		}

		trxDetails = append(trxDetails, model.TrxDetailResp{
			ProductLog: model.ProductLogResp{
				ID:            productID,
				ProductName:   td.ProductLog.ProductName,
				Slug:          td.ProductLog.Slug,
				ResellerPrice: td.ProductLog.ResellerPrice,
				ConsumerPrice: td.ProductLog.ConsumerPrice,
				Description:   td.ProductLog.Description,
				Shop: model.ShopInfo{
					ShopName: td.ProductLog.Shop.ShopName,
					PhotoURL: td.ProductLog.Shop.PhotoURL,
				},
				Category: model.CategoryResp{
					ID:           td.ProductLog.Category.ID,
					CategoryName: td.ProductLog.Category.CategoryName,
				},
				Images: images,
			},
			Shop: model.ShopInfo{
				ShopName: td.ProductLog.Shop.ShopName,
				PhotoURL: td.ProductLog.Shop.PhotoURL,
			},
			Quantity:   td.Quantity,
			TotalPrice: td.TotalPrice,
		})
	}

	res = model.TrxResp{
		ID:            trxResRepo.ID,
		TotalPrice:    trxResRepo.TotalPrice,
		InvoiceCode:   trxResRepo.InvoiceCode,
		PaymentMethod: trxResRepo.PaymentMethod,
		Address: model.AddressResp{
			ID:            addressResRepo.ID,
			AddressTitle:  addressResRepo.AddressTitle,
			RecipientName: addressResRepo.RecipientName,
			PhoneNumber:   addressResRepo.PhoneNumber,
			FullAddress:   addressResRepo.FullAddress,
		},
		TrxDetail: trxDetails,
	}

	return res, nil
}

func (alc *TrxUseCaseImpl) GetAllTrx(ctx context.Context, userID string, params model.TrxFilter) (res model.FilteredData, err *helper.ErrorStruct) {
	var transactions []model.TrxResp

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

	trxesResRepo, errRepo := alc.trxRepository.GetAllTrxByUserID(ctx, userID, entity.FilterTrx{
		Limit:  limit,
		Offset: offset,
		Search: params.Search,
	})
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

	for _, v := range trxesResRepo {
		addressID := fmt.Sprintf("%d", v.AddressID)
		addressResRepo, errRepo := alc.addressesRepository.GetAddressByID(ctx, addressID)
		if !errors.Is(errRepo, gorm.ErrRecordNotFound) && errRepo != nil {
			helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetAddressByID: %s", errRepo.Error()), errRepo)
			return res, &helper.ErrorStruct{
				Code: fiber.StatusBadRequest,
				Err:  errRepo,
			}
		}

		var trxDetails []model.TrxDetailResp
		for _, td := range v.TrxDetails {

			productID := td.ProductLog.ProductID
			imgResRepo, errRepo := alc.productImagesRepository.GetImagesByProductID(ctx, productID)
			if errRepo != nil {
				if errors.Is(errRepo, gorm.ErrRecordNotFound) {
					return res, &helper.ErrorStruct{
						Code: fiber.StatusNotFound,
						Err:  errors.New("foto tidak ditemukan"),
					}
				}

				helper.Logger(helper.LoggerLevelError, fmt.Sprintf("Error at GetImagesByProductID: %s", errRepo.Error()), errRepo)
				return res, &helper.ErrorStruct{
					Code: fiber.StatusBadRequest,
					Err:  errRepo,
				}
			}

			var images []model.ProductImageResp
			for _, img := range imgResRepo {
				images = append(images, model.ProductImageResp{
					ID:        img.ID,
					ProductID: img.ProductID,
					ImageURL:  img.PhotoURL,
				})
			}

			trxDetails = append(trxDetails, model.TrxDetailResp{
				ProductLog: model.ProductLogResp{
					ID:            productID,
					ProductName:   td.ProductLog.ProductName,
					Slug:          td.ProductLog.Slug,
					ResellerPrice: td.ProductLog.ResellerPrice,
					ConsumerPrice: td.ProductLog.ConsumerPrice,
					Description:   td.ProductLog.Description,
					Shop: model.ShopInfo{
						ShopName: td.ProductLog.Shop.ShopName,
						PhotoURL: td.ProductLog.Shop.PhotoURL,
					},
					Category: model.CategoryResp{
						ID:           td.ProductLog.Category.ID,
						CategoryName: td.ProductLog.Category.CategoryName,
					},
					Images: images,
				},
				Shop: model.ShopInfo{
					ShopName: td.ProductLog.Shop.ShopName,
					PhotoURL: td.ProductLog.Shop.PhotoURL,
				},
				Quantity:   td.Quantity,
				TotalPrice: td.TotalPrice,
			})
		}

		transactions = append(transactions, model.TrxResp{
			ID:            v.ID,
			TotalPrice:    v.TotalPrice,
			InvoiceCode:   v.InvoiceCode,
			PaymentMethod: v.PaymentMethod,
			Address: model.AddressResp{
				ID:            addressResRepo.ID,
				AddressTitle:  addressResRepo.AddressTitle,
				RecipientName: addressResRepo.RecipientName,
				PhoneNumber:   addressResRepo.PhoneNumber,
				FullAddress:   addressResRepo.FullAddress,
			},
			TrxDetail: trxDetails,
		})
	}

	res = model.FilteredData{
		Data:  transactions,
		Page:  params.Page,
		Limit: params.Limit,
	}

	return res, err
}
