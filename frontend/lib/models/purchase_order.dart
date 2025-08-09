import 'package:json_annotation/json_annotation.dart';

part 'purchase_order.g.dart';

@JsonSerializable()
class PurchaseOrder {
  final int id;
  @JsonKey(name: 'po_number')
  final String poNumber;
  @JsonKey(name: 'requisition_id')
  final int requisitionId;
  @JsonKey(name: 'vendor_id')
  final int vendorId;
  @JsonKey(name: 'order_date')
  final DateTime orderDate;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  PurchaseOrder({
    required this.id,
    required this.poNumber,
    required this.requisitionId,
    required this.vendorId,
    required this.orderDate,
    required this.createdAt,
  });

  factory PurchaseOrder.fromJson(Map<String, dynamic> json) => _$PurchaseOrderFromJson(json);
  Map<String, dynamic> toJson() => _$PurchaseOrderToJson(this);
}
