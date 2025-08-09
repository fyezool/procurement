// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'navigation_item.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

NavigationItem _$NavigationItemFromJson(Map<String, dynamic> json) =>
    NavigationItem(
      title: json['title'] as String,
      path: json['path'] as String,
      icon: json['icon'] as String,
      subItems: (json['subItems'] as List<dynamic>?)
              ?.map(
                  (e) => NavigationSubItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          const [],
    );

Map<String, dynamic> _$NavigationItemToJson(NavigationItem instance) =>
    <String, dynamic>{
      'title': instance.title,
      'path': instance.path,
      'icon': instance.icon,
      'subItems': instance.subItems,
    };

NavigationSubItem _$NavigationSubItemFromJson(Map<String, dynamic> json) =>
    NavigationSubItem(
      title: json['title'] as String,
      path: json['path'] as String,
    );

Map<String, dynamic> _$NavigationSubItemToJson(NavigationSubItem instance) =>
    <String, dynamic>{
      'title': instance.title,
      'path': instance.path,
    };
