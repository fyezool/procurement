class Contract {
  final String? id;
  final String procurementId;
  final String supplierId;
  final List<ContractItem> items;
  final double totalCost;

  Contract({
    this.id,
    required this.procurementId,
    required this.supplierId,
    required this.items,
    required this.totalCost,
  });

  factory Contract.fromJson(Map<String, dynamic> json) {
    var itemsJson = json['items'] as List;
    List<ContractItem> items = itemsJson.map((i) => ContractItem.fromJson(i)).toList();

    return Contract(
      id: json['_id'],
      procurementId: json['procurement_id'],
      supplierId: json['supplier_id'],
      items: items,
      totalCost: json['total_cost'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'procurement_id': procurementId,
      'supplier_id': supplierId,
      'items': items.map((item) => item.toJson()).toList(),
      'total_cost': totalCost,
    };
  }
}

class ContractItem {
  final String itemId;
  final int quantity;
  final double unitPrice;

  ContractItem({
    required this.itemId,
    required this.quantity,
    required this.unitPrice,
  });

  factory ContractItem.fromJson(Map<String, dynamic> json) {
    return ContractItem(
      itemId: json['item_id'],
      quantity: json['quantity'],
      unitPrice: json['unit_price'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'item_id': itemId,
      'quantity': quantity,
      'unit_price': unitPrice,
    };
  }
}