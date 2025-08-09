import 'package:flutter/material.dart';

class PendingRequisitionsScreen extends StatelessWidget {
  const PendingRequisitionsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Pending Requisitions'),
      ),
      body: const Center(
        child: Text('Pending Requisitions - Under Development'),
      ),
    );
  }
}
