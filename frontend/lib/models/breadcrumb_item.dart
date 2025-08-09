import 'package:json_annotation/json_annotation.dart';

part 'breadcrumb_item.g.dart';

@JsonSerializable()
class BreadcrumbItem {
  final String title;
  final String path;

  BreadcrumbItem({
    required this.title,
    required this.path,
  });

  factory BreadcrumbItem.fromJson(Map<String, dynamic> json) => _$BreadcrumbItemFromJson(json);
  Map<String, dynamic> toJson() => _$BreadcrumbItemToJson(this);
}
