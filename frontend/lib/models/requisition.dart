import 'package:json_annotation/json_annotation.dart';

part 'requisition.g.dart';

@JsonSerializable()
class Requisition {
  final int id;
  @JsonKey(name: 'requester_id')
  final int requesterId;
  @JsonKey(name: 'vendor_id')
  final int? vendorId;
  @JsonKey(name: 'item_description')
  final String itemDescription;
  final int quantity;
  @JsonKey(name: 'estimated_price')
  final double estimatedPrice;
  @JsonKey(name: 'total_price')
  final double totalPrice;
  final String? justification;
  final String status;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  Requisition({
    required this.id,
    required this.requesterId,
    this.vendorId,
    required this.itemDescription,
    required this.quantity,
    required this.estimatedPrice,
    required this.totalPrice,
    this.justification,
    required this.status,
    required this.createdAt,
  });

  factory Requisition.fromJson(Map<String, dynamic> json) => _$RequisitionFromJson(json);
  Map<String, dynamic> toJson() => _$RequisitionToJson(this);
}
