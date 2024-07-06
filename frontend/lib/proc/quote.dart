class Quote {
  final String? id;
  final String procurementId;
  final String supplierId;
  final List<QuoteItem> items;
  final double totalCost;

  Quote({
    this.id,
    required this.procurementId,
    required this.supplierId,
    required this.items,
    required this.totalCost,
  });

  factory Quote.fromJson(Map<String, dynamic> json) {
    var itemsJson = json['items'] as List;
    List<QuoteItem> items = itemsJson.map((i) => QuoteItem.fromJson(i)).toList();

    return Quote(
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

class QuoteItem {
  final String itemId;
  final int quantity;
  final double unitPrice;

  QuoteItem({
    required this.itemId,
    required this.quantity,
    required this.unitPrice,
  });

  factory QuoteItem.fromJson(Map<String, dynamic> json) {
    return QuoteItem(
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