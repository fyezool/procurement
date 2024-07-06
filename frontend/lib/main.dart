import 'package:flutter/material.dart';
import 'screens/procurement_list_screen.dart';
import 'screens/supplier_list_screen.dart';
import 'screens/item_list_screen.dart';
import 'screens/quote_list_screen.dart';
import 'screens/contract_list_screen.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Procurement App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: HomeScreen(),
    );
  }
}

class HomeScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Procurement App'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ProcurementListScreen()),
                );
              },
              child: Text('Procurements'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => SupplierListScreen()),
                );
              },
              child: Text('Suppliers'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ItemListScreen()),
                );
              },
              child: Text('Items'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => QuoteListScreen()),
                );
              },
              child: Text('Quotes'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ContractListScreen()),
                );
              },
              child: Text('Contracts'),
            ),
          ],
        ),
      ),
    );
  }
}