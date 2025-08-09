import 'package:json_annotation/json_annotation.dart';

part 'navigation_item.g.dart';

@JsonSerializable()
class NavigationItem {
  final String title;
  final String path;
  final String icon;
  final List<NavigationSubItem> subItems;

  NavigationItem({
    required this.title,
    required this.path,
    required this.icon,
    this.subItems = const [],
  });

  factory NavigationItem.fromJson(Map<String, dynamic> json) => _$NavigationItemFromJson(json);
  Map<String, dynamic> toJson() => _$NavigationItemToJson(this);
}

@JsonSerializable()
class NavigationSubItem {
  final String title;
  final String path;

  NavigationSubItem({
    required this.title,
    required this.path,
  });

  factory NavigationSubItem.fromJson(Map<String, dynamic> json) => _$NavigationSubItemFromJson(json);
  Map<String, dynamic> toJson() => _$NavigationSubItemToJson(this);
}
