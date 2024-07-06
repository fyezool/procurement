import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../proc/supplier.dart';

class SupplierListScreen extends StatefulWidget {
  @override
  _SupplierListScreenState createState() => _SupplierListScreenState();
}

class _SupplierListScreenState extends State<SupplierListScreen> {
  List<Supplier> suppliers = [];

  @override
  void initState() {
    super.initState();
    fetchSuppliers();
  }

  Future<void> fetchSuppliers() async {
    final response = await http.get(Uri.parse('http://localhost:8080/suppliers'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        suppliers = data.map((supplier) => Supplier.fromJson(supplier)).toList();
      });
    }
  }

  Future<void> createSupplier(Supplier supplier) async {
    final response = await http.post(
      Uri.parse('http://localhost:8080/suppliers'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(supplier.toJson()),
    );
    if (response.statusCode == 200) {
      fetchSuppliers();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Suppliers'),
      ),
      body: ListView.builder(
        itemCount: suppliers.length,
        itemBuilder: (context, index) {
          return ListTile(
            title: Text(suppliers[index].name),
            subtitle: Text('Contact: ${suppliers[index].contactInfo}'),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
            context: context,
            builder: (context) {
              return AlertDialog(
                title: Text('Create Supplier'),
                content: SupplierForm(
                  onSubmit: (supplier) {
                    createSupplier(supplier);
                    Navigator.of(context).pop();
                  },
                ),
              );
            },
          );
        },
        child: Icon(Icons.add),
      ),
    );
  }
}

class SupplierForm extends StatefulWidget {
  final Function(Supplier) onSubmit;

  SupplierForm({required this.onSubmit});

  @override
  _SupplierFormState createState() => _SupplierFormState();
}

class _SupplierFormState extends State<SupplierForm> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _contactInfoController = TextEditingController();
  final _addressController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextFormField(
            controller: _nameController,
            decoration: InputDecoration(labelText: 'Name'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a name';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _contactInfoController,
            decoration: InputDecoration(labelText: 'Contact Info'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter contact information';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _addressController,
            decoration: InputDecoration(labelText: 'Address'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter an address';
              }
              return null;
            },
          ),
          ElevatedButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                widget.onSubmit(
                  Supplier(
                    name: _nameController.text,
                    contactInfo: _contactInfoController.text,
                    address: _addressController.text,
                  ),
                );
              }
            },
            child: Text('Submit'),
          ),
        ],
      ),
    );
  }
}