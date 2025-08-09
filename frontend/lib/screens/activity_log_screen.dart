import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../models/activity_log.dart';
import '../services/api_service.dart';
import '../widgets/empty_state_widget.dart';

class ActivityLogScreen extends StatefulWidget {
  const ActivityLogScreen({Key? key}) : super(key: key);

  @override
  _ActivityLogScreenState createState() => _ActivityLogScreenState();
}

class _ActivityLogScreenState extends State<ActivityLogScreen> {
  late Future<List<ActivityLog>> _logsFuture;
  final ApiService _apiService = ApiService();

  @override
  void initState() {
    super.initState();
    _logsFuture = _apiService.getActivityLogs();
  }

  void _refreshLogs() {
    setState(() {
      _logsFuture = _apiService.getActivityLogs();
    });
  }

  IconData _getIconForStatus(String status) {
    return status == 'SUCCESS' ? Icons.check_circle : Icons.error;
  }

  Color _getColorForStatus(String status) {
    return status == 'SUCCESS' ? Colors.green : Colors.red;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Activity Log'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshLogs,
          ),
        ],
      ),
      body: FutureBuilder<List<ActivityLog>>(
        future: _logsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return EmptyStateWidget(
              message: 'Failed to load activity logs: ${snapshot.error}',
              icon: Icons.error_outline,
              onRetry: _refreshLogs,
            );
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return EmptyStateWidget(
              message: 'No activity has been logged yet.',
              icon: Icons.history_toggle_off,
              onRetry: _refreshLogs,
            );
          }

          final logs = snapshot.data!;
          return ListView.builder(
            itemCount: logs.length,
            itemBuilder: (context, index) {
              final log = logs[index];
              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
                child: ListTile(
                  leading: Icon(
                    _getIconForStatus(log.status),
                    color: _getColorForStatus(log.status),
                  ),
                  title: Text(log.action),
                  subtitle: Text(
                    'User: ${log.userId ?? 'System'} | Target: ${log.targetType ?? 'N/A'} #${log.targetId ?? ''}\n'
                    'Details: ${log.details ?? 'No details'}'
                  ),
                  trailing: Text(
                    DateFormat('yyyy-MM-dd HH:mm').format(log.createdAt.toLocal()),
                  ),
                  isThreeLine: true,
                ),
              );
            },
          );
        },
      ),
    );
  }
}
