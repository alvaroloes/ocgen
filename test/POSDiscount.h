//
//  POSDiscount.h
//  PointOfSaleBlock
//
//  Created by Juan Fernández Sagasti on 23/9/15.
//  Copyright © 2015 eBay, Inc. All rights reserved.
//

typedef NS_ENUM(NSInteger, POSDiscountType)
{
    POSDiscountTypeSpecificPrice,
    POSDiscountTypeAmountOff,
    POSDiscountTypePercentageOff
};


@interface POSDiscount : NSObject <POSModelProtocol> OCGEN_AUTO

@property (nonatomic, strong) NSNumber *discountAdjustedAmount;
@property (nonatomic, strong) NSNumber *discountAmount;
@property (nonatomic, strong) NSNumber *discountBasisAmount;
@property (nonatomic, copy) NSString *discountCode;
@property (nonatomic, copy) NSString *discountDescription;
@property (nonatomic, copy) NSString *discountID;
@property (nonatomic, assign) BOOL discountIsApplied;
@property (nonatomic, assign) BOOL discountIsTaxable;
@property (nonatomic, assign) BOOL discountIsStackable;
@property (nonatomic, assign) POSDiscountType discountType;
@property (nonatomic, strong) NSNumber *discountValue;

@end
