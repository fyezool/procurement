class Item {
  final String? id;
  final String name;
  final String description;
  final int quantity;
  final double unitPrice;

  Item({
    this.id,
    required this.name,
    required this.description,
    required this.quantity,
    required this.unitPrice,
  });

  factory Item.fromJson(Map<String, dynamic> json) {
    return Item(
      id: json['_id'],
      name: json['name'],
      description: json['description'],
      quantity: json['quantity'],
      unitPrice: json['unit_price'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'description': description,
      'quantity': quantity,
      'unit_price': unitPrice,
    };
  }
}