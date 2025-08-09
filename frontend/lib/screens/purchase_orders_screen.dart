import 'package:flutter/material.dart';
import '../models/purchase_order.dart';
import '../services/api_service.dart';

class PurchaseOrdersScreen extends StatefulWidget {
  const PurchaseOrdersScreen({Key? key}) : super(key: key);

  @override
  _PurchaseOrdersScreenState createState() =>
      _PurchaseOrdersScreenState();
}

class _PurchaseOrdersScreenState extends State<PurchaseOrdersScreen> {
  late Future<List<PurchaseOrder>> _posFuture;

  @override
  void initState() {
    super.initState();
    // For now, this screen will show all POs, assuming admin access.
    // A role check could be added here to fetch user-specific POs later.
    _posFuture = ApiService().getAllPurchaseOrders();
  }

  void _refreshPOs() {
    setState(() {
      _posFuture = ApiService().getAllPurchaseOrders();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Purchase Orders'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshPOs,
          ),
        ],
      ),
      body: FutureBuilder<List<PurchaseOrder>>(
        future: _posFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('No purchase orders found.'));
          }

          final pos = snapshot.data!;
          return SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: DataTable(
                columns: const [
                  DataColumn(label: Text('ID')),
                  DataColumn(label: Text('PO Number')),
                  DataColumn(label: Text('Requisition ID')),
                  DataColumn(label: Text('Vendor ID')),
                  DataColumn(label: Text('Order Date')),
                ],
                rows: pos.map((po) {
                  return DataRow(cells: [
                    DataCell(Text(po.id.toString())),
                    DataCell(Text(po.poNumber)),
                    DataCell(Text(po.requisitionId.toString())),
                    DataCell(Text(po.vendorId.toString())),
                    DataCell(Text(po.orderDate.toLocal().toString().split(' ')[0])),
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
