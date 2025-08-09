import 'package:flutter/material.dart';
import '../models/vendor.dart';
import '../services/api_service.dart';

class VendorsScreen extends StatefulWidget {
  const VendorsScreen({Key? key}) : super(key: key);

  @override
  _VendorsScreenState createState() => _VendorsScreenState();
}

class _VendorsScreenState extends State<VendorsScreen> {
  late Future<List<Vendor>> _vendorsFuture;

  @override
  void initState() {
    super.initState();
    _vendorsFuture = ApiService().getVendors();
  }

  void _refreshVendors() {
    setState(() {
      _vendorsFuture = ApiService().getVendors();
    });
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
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('No vendors found.'));
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
                ],
                rows: vendors.map((vendor) {
                  return DataRow(cells: [
                    DataCell(Text(vendor.id.toString())),
                    DataCell(Text(vendor.name)),
                    DataCell(Text(vendor.contactPerson ?? 'N/A')),
                    DataCell(Text(vendor.email ?? 'N/A')),
                    DataCell(Text(vendor.phone ?? 'N/A')),
                  ]);
                }).toList(),
              ),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // TODO: Navigate to create vendor screen
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Create Vendor screen not implemented yet.')),
          );
        },
        child: const Icon(Icons.add),
        tooltip: 'Add Vendor',
      ),
    );
  }
}
