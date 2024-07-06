class Supplier {
  final String? id;
  final String name;
  final String contactInfo;
  final String address;

  Supplier({
    this.id,
    required this.name,
    required this.contactInfo,
    required this.address,
  });

  factory Supplier.fromJson(Map<String, dynamic> json) {
    return Supplier(
      id: json['_id'],
      name: json['name'],
      contactInfo: json['contact_info'],
      address: json['address'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'contact_info': contactInfo,
      'address': address,
    };
  }
}