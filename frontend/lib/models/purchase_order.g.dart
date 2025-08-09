// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'purchase_order.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

PurchaseOrder _$PurchaseOrderFromJson(Map<String, dynamic> json) =>
    PurchaseOrder(
      id: (json['id'] as num).toInt(),
      poNumber: json['po_number'] as String,
      requisitionId: (json['requisition_id'] as num).toInt(),
      vendorId: (json['vendor_id'] as num).toInt(),
      orderDate: DateTime.parse(json['order_date'] as String),
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$PurchaseOrderToJson(PurchaseOrder instance) =>
    <String, dynamic>{
      'id': instance.id,
      'po_number': instance.poNumber,
      'requisition_id': instance.requisitionId,
      'vendor_id': instance.vendorId,
      'order_date': instance.orderDate.toIso8601String(),
      'created_at': instance.createdAt.toIso8601String(),
    };
