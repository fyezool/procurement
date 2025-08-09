import 'package:flutter/material.dart';
import '../../models/requisition.dart';
import '../../services/api_service.dart';

class AllRequisitionsScreen extends StatefulWidget {
  const AllRequisitionsScreen({Key? key}) : super(key: key);

  @override
  _AllRequisitionsScreenState createState() => _AllRequisitionsScreenState();
}

class _AllRequisitionsScreenState extends State<AllRequisitionsScreen> {
  late Future<List<Requisition>> _requisitionsFuture;

  @override
  void initState() {
    super.initState();
    _requisitionsFuture = ApiService().getAllRequisitions();
  }

  void _refreshRequisitions() {
    setState(() {
      _requisitionsFuture = ApiService().getAllRequisitions();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('All Requisitions'),
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
            return const Center(child: Text('No requisitions found.'));
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
                  DataColumn(label: Text('Status')),
                  DataColumn(label: Text('Total Price')),
                ],
                rows: requisitions.map((req) {
                  return DataRow(cells: [
                    DataCell(Text(req.id.toString())),
                    DataCell(Text(req.requesterId.toString())),
                    DataCell(Text(req.itemDescription)),
                    DataCell(Text(req.status)),
                    DataCell(Text('\$${req.totalPrice.toStringAsFixed(2)}')),
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
