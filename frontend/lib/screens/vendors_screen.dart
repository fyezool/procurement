import 'package:flutter/material.dart';

class VendorsScreen extends StatelessWidget {
  const VendorsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Vendors'),
      ),
      body: const Center(
        child: Text('Vendor Management'),
      ),
    );
  }
}
