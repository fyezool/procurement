import 'package:flutter/material.dart';

class ProcurementMainScreen extends StatelessWidget {
  const ProcurementMainScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Procurement'),
      ),
      body: const Center(
        child: Text('Procurement Main Screen'),
      ),
    );
  }
}
