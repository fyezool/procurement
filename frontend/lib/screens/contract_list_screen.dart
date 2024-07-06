import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../proc/contract.dart';

class ContractListScreen extends StatefulWidget {
  @override
  _ContractListScreenState createState() => _ContractListScreenState();
}

class _ContractListScreenState extends State<ContractListScreen> {
  List<Contract> contracts = [];

  @override
  void initState() {
    super.initState();
    fetchContracts();
  }

  Future<void> fetchContracts() async {
    final response = await http.get(Uri.parse('http://localhost:8080/contracts'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        contracts = data.map((contract) => Contract.fromJson(contract)).toList();
      });
    }
  }

  Future<void> createContract(Contract contract) async {
    final response = await http.post(
      Uri.parse('http://localhost:8080/contracts'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(contract.toJson()),
    );
    if (response.statusCode == 200) {
      fetchContracts();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Contracts'),
      ),
      body: ListView.builder(
        itemCount: contracts.length,
        itemBuilder: (context, index) {
          return ListTile(
            title: Text('Contract ${contracts[index].id}'),
            subtitle: Text('Total Cost: \$${contracts[index].totalCost}'),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
            context: context,
            builder: (context) {
              return AlertDialog(
                title: Text('Create Contract'),
                content: ContractForm(
                  onSubmit: (contract) {
                    createContract(contract);
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

class ContractForm extends StatefulWidget {
  final Function(Contract) onSubmit;

  ContractForm({required this.onSubmit});

  @override
  _ContractFormState createState() => _ContractFormState();
}

class _ContractFormState extends State<ContractForm> {
  final _formKey = GlobalKey<FormState>();
  final _procurementIdController = TextEditingController();
  final _supplierIdController = TextEditingController();
  final _totalCostController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextFormField(
            controller: _procurementIdController,
            decoration: InputDecoration(labelText: 'Procurement ID'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a procurement ID';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _supplierIdController,
            decoration: InputDecoration(labelText: 'Supplier ID'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a supplier ID';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _totalCostController,
            decoration: InputDecoration(labelText: 'Total Cost'),
            keyboardType: TextInputType.number,
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a total cost';
              }
              return null;
            },
          ),
          ElevatedButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                widget.onSubmit(
                  Contract(
                    procurementId: _procurementIdController.text,
                    supplierId: _supplierIdController.text,
                    items: [], // You can add items here or handle them separately
                    totalCost: double.parse(_totalCostController.text),
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