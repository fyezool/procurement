import 'package:json_annotation/json_annotation.dart';

part 'activity_log.g.dart';

@JsonSerializable()
class ActivityLog {
  final int id;
  @JsonKey(name: 'user_id')
  final int? userId;
  final String action;
  @JsonKey(name: 'target_type')
  final String? targetType;
  @JsonKey(name: 'target_id')
  final int? targetId;
  final String status;
  final String? details;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  ActivityLog({
    required this.id,
    this.userId,
    required this.action,
    this.targetType,
    this.targetId,
    required this.status,
    this.details,
    required this.createdAt,
  });

  factory ActivityLog.fromJson(Map<String, dynamic> json) => _$ActivityLogFromJson(json);
  Map<String, dynamic> toJson() => _$ActivityLogToJson(this);
}
