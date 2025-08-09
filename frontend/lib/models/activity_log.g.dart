// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'activity_log.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ActivityLog _$ActivityLogFromJson(Map<String, dynamic> json) => ActivityLog(
      id: (json['id'] as num).toInt(),
      userId: (json['user_id'] as num?)?.toInt(),
      action: json['action'] as String,
      targetType: json['target_type'] as String?,
      targetId: (json['target_id'] as num?)?.toInt(),
      status: json['status'] as String,
      details: json['details'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$ActivityLogToJson(ActivityLog instance) =>
    <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'action': instance.action,
      'target_type': instance.targetType,
      'target_id': instance.targetId,
      'status': instance.status,
      'details': instance.details,
      'created_at': instance.createdAt.toIso8601String(),
    };
