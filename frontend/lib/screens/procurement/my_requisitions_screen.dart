import 'package:flutter/material.dart';
import '../../models/requisition.dart';
import '../../services/api_service.dart';
import '../../widgets/edit_requisition_dialog.dart';

class MyRequisitionsScreen extends StatefulWidget {
  const MyRequisitionsScreen({Key? key}) : super(key: key);

  @override
  _MyRequisitionsScreenState createState() => _MyRequisitionsScreenState();
}

class _MyRequisitionsScreenState extends State<MyRequisitionsScreen> {
  late Future<List<Requisition>> _requisitionsFuture;
  final ApiService _apiService = ApiService();

  @override
  void initState() {
    super.initState();
    _requisitionsFuture = _apiService.getMyRequisitions();
  }

  void _refreshRequisitions() {
    setState(() {
      _requisitionsFuture = _apiService.getMyRequisitions();
    });
  }

  void _showEditRequisitionDialog(Requisition requisition) {
    showDialog(
      context: context,
      builder: (context) {
        return EditRequisitionDialog(
          requisition: requisition,
          onSave: (updatedData) async {
            try {
              await _apiService.updateRequisition(requisition.id, updatedData);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Requisition updated successfully'),
                  backgroundColor: Colors.green,
                ),
              );
              _refreshRequisitions();
            } catch (e) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('Failed to update requisition: $e'),
                  backgroundColor: Colors.red,
                ),
              );
            }
          },
        );
      },
    );
  }

  void _showDeleteConfirmationDialog(Requisition requisition) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Delete Requisition'),
          content: Text('Are you sure you want to delete this requisition: "${requisition.itemDescription}"?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            TextButton(
              style: TextButton.styleFrom(
                foregroundColor: Colors.white,
                backgroundColor: Colors.red,
              ),
              onPressed: () async {
                try {
                  await _apiService.deleteRequisition(requisition.id);
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Requisition deleted successfully'),
                      backgroundColor: Colors.green,
                    ),
                  );
                  _refreshRequisitions();
                } catch (e) {
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Failed to delete requisition: $e'),
                      backgroundColor: Colors.red,
                    ),
                  );
                }
              },
              child: const Text('Delete'),
            ),
          ],
        );
      },
    );
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Requisitions'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshRequisitions,
          ),
        ],
      ),
      body: FutureBuilder<List<Requisition>>(
        future: _requisitionsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('You have not created any requisitions.'));
          }

          final requisitions = snapshot.data!;
          return SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: DataTable(
                columns: const [
                  DataColumn(label: Text('ID')),
                  DataColumn(label: Text('Description')),
                  DataColumn(label: Text('Status')),
                  DataColumn(label: Text('Total Price')),
                  DataColumn(label: Text('Actions')),
                ],
                rows: requisitions.map((req) {
                  final bool isPending = req.status == 'Pending';
                  return DataRow(cells: [
                    DataCell(Text(req.id.toString())),
                    DataCell(Text(req.itemDescription)),
                    DataCell(Text(req.status)),
                    DataCell(Text('\$${req.totalPrice.toStringAsFixed(2)}')),
                    DataCell(Row(
                      children: [
                        if (isPending)
                          IconButton(
                            icon: const Icon(Icons.edit),
                            onPressed: () => _showEditRequisitionDialog(req),
                          ),
                        if (isPending)
                          IconButton(
                            icon: const Icon(Icons.delete, color: Colors.red),
                            onPressed: () => _showDeleteConfirmationDialog(req),
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
