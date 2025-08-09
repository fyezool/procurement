import 'package:flutter/material.dart';

class PurchaseOrdersScreen extends StatelessWidget {
  const PurchaseOrdersScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Purchase Orders'),
      ),
      body: const Center(
        child: Text('Purchase Order Management'),
      ),
    );
  }
}
