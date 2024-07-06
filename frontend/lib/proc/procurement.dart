class Procurement {
  final String? id;
  final String title;
  final String description;
  final String status;
  final DateTime dateCreated;
  final DateTime dateUpdated;

  Procurement({
    this.id,
    required this.title,
    required this.description,
    required this.status,
    required this.dateCreated,
    required this.dateUpdated,
  });

  factory Procurement.fromJson(Map<String, dynamic> json) {
    return Procurement(
      id: json['_id'],
      title: json['title'],
      description: json['description'],
      status: json['status'],
      dateCreated: DateTime.parse(json['date_created']),
      dateUpdated: DateTime.parse(json['date_updated']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id' : id,
      'title': title,
      'description': description,
      'status': status,
      'date_created': dateCreated.toIso8601String(),
      'date_updated': dateUpdated.toIso8601String(),
    };
  }
}