import 'package:flutter/material.dart';
import '../../models/requisition.dart';
import '../../services/api_service.dart';
import '../../widgets/empty_state_widget.dart';

class ApprovalsScreen extends StatefulWidget {
  const ApprovalsScreen({Key? key}) : super(key: key);

  @override
  _ApprovalsScreenState createState() => _ApprovalsScreenState();
}

class _ApprovalsScreenState extends State<ApprovalsScreen> {
  late Future<List<Requisition>> _pendingRequisitionsFuture;
  final ApiService _apiService = ApiService();

  @override
  void initState() {
    super.initState();
    _pendingRequisitionsFuture = _apiService.getPendingRequisitions();
  }

  void _refreshPendingRequisitions() {
    setState(() {
      _pendingRequisitionsFuture = _apiService.getPendingRequisitions();
    });
  }

  void _approve(int id) async {
    try {
      await _apiService.approveRequisition(id);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Requisition Approved'), backgroundColor: Colors.green),
      );
      _refreshPendingRequisitions();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to approve: $e'), backgroundColor: Colors.red),
      );
    }
  }

  void _reject(int id) async {
    try {
      await _apiService.rejectRequisition(id);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Requisition Rejected'), backgroundColor: Colors.orange),
      );
      _refreshPendingRequisitions();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to reject: $e'), backgroundColor: Colors.red),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Pending Approvals'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshPendingRequisitions,
          ),
        ],
      ),
      body: FutureBuilder<List<Requisition>>(
        future: _pendingRequisitionsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return EmptyStateWidget(
              message: 'Failed to load pending requisitions: ${snapshot.error}',
              icon: Icons.error_outline,
              onRetry: _refreshPendingRequisitions,
            );
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return EmptyStateWidget(
              message: 'There are no requisitions waiting for approval.',
              icon: Icons.check_circle_outline,
              onRetry: _refreshPendingRequisitions,
            );
          }

          final requisitions = snapshot.data!;
          return SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: DataTable(
                columns: const [
                  DataColumn(label: Text('ID')),
                  DataColumn(label: Text('Requester ID')),
                  DataColumn(label: Text('Description')),
                  DataColumn(label: Text('Total Price')),
                  DataColumn(label: Text('Actions')),
                ],
                rows: requisitions.map((req) {
                  return DataRow(cells: [
                    DataCell(Text(req.id.toString())),
                    DataCell(Text(req.requesterId.toString())),
                    DataCell(Text(req.itemDescription)),
                    DataCell(Text('\$${req.totalPrice.toStringAsFixed(2)}')),
                    DataCell(Row(
                      children: [
                        ElevatedButton(
                          onPressed: () => _approve(req.id),
                          child: const Text('Approve'),
                          style: ElevatedButton.styleFrom(backgroundColor: Colors.green),
                        ),
                        const SizedBox(width: 8),
                        ElevatedButton(
                          onPressed: () => _reject(req.id),
                          child: const Text('Reject'),
                          style: ElevatedButton.styleFrom(backgroundColor: Colors.orange),
                        ),
                      ],
                    )),
                  ]);
                }).toList(),
              ),
            ),
          );
        },
      ),
    );
  }
}
