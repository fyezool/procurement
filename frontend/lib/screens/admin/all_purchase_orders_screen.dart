import 'package:flutter/material.dart';
import '../../models/purchase_order.dart';
import '../../services/api_service.dart';
import '../../widgets/empty_state_widget.dart';

class AllPurchaseOrdersScreen extends StatefulWidget {
  const AllPurchaseOrdersScreen({Key? key}) : super(key: key);

  @override
  _AllPurchaseOrdersScreenState createState() =>
      _AllPurchaseOrdersScreenState();
}

class _AllPurchaseOrdersScreenState extends State<AllPurchaseOrdersScreen> {
  late Future<List<PurchaseOrder>> _posFuture;

  @override
  void initState() {
    super.initState();
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
        title: const Text('All Purchase Orders'),
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
            return EmptyStateWidget(
              message: 'Failed to load purchase orders: ${snapshot.error}',
              icon: Icons.error_outline,
              onRetry: _refreshPOs,
            );
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return EmptyStateWidget(
              message: 'No purchase orders found in the system.',
              icon: Icons.receipt_long_outlined,
              onRetry: _refreshPOs,
            );
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
