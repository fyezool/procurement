import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../models/vendor.dart';
import '../services/api_service.dart';
import '../widgets/edit_vendor_dialog.dart';
import '../widgets/empty_state_widget.dart';

class VendorsScreen extends StatefulWidget {
  const VendorsScreen({Key? key}) : super(key: key);

  @override
  _VendorsScreenState createState() => _VendorsScreenState();
}

class _VendorsScreenState extends State<VendorsScreen> {
  late Future<List<Vendor>> _vendorsFuture;
  final ApiService _apiService = ApiService();

  @override
  void initState() {
    super.initState();
    _vendorsFuture = _apiService.getVendors();
  }

  void _refreshVendors() {
    setState(() {
      _vendorsFuture = _apiService.getVendors();
    });
  }

  void _showEditVendorDialog(Vendor vendor) {
    showDialog(
      context: context,
      builder: (context) {
        return EditVendorDialog(
          vendor: vendor,
          onSave: (updatedVendor) async {
            try {
              await _apiService.updateVendor(vendor.id, updatedVendor.toJson());
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Vendor updated successfully'),
                  backgroundColor: Colors.green,
                ),
              );
              _refreshVendors();
            } catch (e) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('Failed to update vendor: $e'),
                  backgroundColor: Colors.red,
                ),
              );
            }
          },
        );
      },
    );
  }

  void _showDeleteConfirmationDialog(Vendor vendor) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Delete Vendor'),
          content: Text('Are you sure you want to delete ${vendor.name}?'),
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
                  await _apiService.deleteVendor(vendor.id);
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Vendor deleted successfully'),
                      backgroundColor: Colors.green,
                    ),
                  );
                  _refreshVendors();
                } catch (e) {
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Failed to delete vendor: $e'),
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
        title: const Text('Vendors'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshVendors,
          ),
        ],
      ),
      body: FutureBuilder<List<Vendor>>(
        future: _vendorsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return EmptyStateWidget(
              message: 'Failed to load vendors: ${snapshot.error}',
              icon: Icons.error_outline,
              onRetry: _refreshVendors,
            );
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return EmptyStateWidget(
              message: 'No vendors found. Tap the button to add one.',
              icon: Icons.store_mall_directory_outlined,
              onRetry: _refreshVendors,
            );
          }

          final vendors = snapshot.data!;
          return SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: DataTable(
                columns: const [
                  DataColumn(label: Text('ID')),
                  DataColumn(label: Text('Name')),
                  DataColumn(label: Text('Contact Person')),
                  DataColumn(label: Text('Email')),
                  DataColumn(label: Text('Phone')),
                  DataColumn(label: Text('Actions')),
                ],
                rows: vendors.map((vendor) {
                  return DataRow(cells: [
                    DataCell(Text(vendor.id.toString())),
                    DataCell(Text(vendor.name)),
                    DataCell(Text(vendor.contactPerson ?? 'N/A')),
                    DataCell(Text(vendor.email ?? 'N/A')),
                    DataCell(Text(vendor.phone ?? 'N/A')),
                    DataCell(Row(
                      children: [
                        IconButton(
                          icon: const Icon(Icons.edit),
                          onPressed: () => _showEditVendorDialog(vendor),
                        ),
                        IconButton(
                          icon: const Icon(Icons.delete, color: Colors.red),
                          onPressed: () => _showDeleteConfirmationDialog(vendor),
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
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/vendors/create'),
        child: const Icon(Icons.add),
        tooltip: 'Add Vendor',
      ),
    );
  }
}
