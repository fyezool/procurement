// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'requisition.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Requisition _$RequisitionFromJson(Map<String, dynamic> json) => Requisition(
      id: (json['id'] as num).toInt(),
      requesterId: (json['requester_id'] as num).toInt(),
      vendorId: (json['vendor_id'] as num?)?.toInt(),
      itemDescription: json['item_description'] as String,
      quantity: (json['quantity'] as num).toInt(),
      estimatedPrice: (json['estimated_price'] as num).toDouble(),
      totalPrice: (json['total_price'] as num).toDouble(),
      justification: json['justification'] as String?,
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$RequisitionToJson(Requisition instance) =>
    <String, dynamic>{
      'id': instance.id,
      'requester_id': instance.requesterId,
      'vendor_id': instance.vendorId,
      'item_description': instance.itemDescription,
      'quantity': instance.quantity,
      'estimated_price': instance.estimatedPrice,
      'total_price': instance.totalPrice,
      'justification': instance.justification,
      'status': instance.status,
      'created_at': instance.createdAt.toIso8601String(),
    };
