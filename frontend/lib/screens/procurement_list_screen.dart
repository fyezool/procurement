import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../proc/procurement.dart';

class ProcurementListScreen extends StatefulWidget {
  @override
  _ProcurementListScreenState createState() => _ProcurementListScreenState();
}

class _ProcurementListScreenState extends State<ProcurementListScreen> {
  List<Procurement> procurements = [];

  @override
  void initState() {
    super.initState();
    fetchProcurements();
  }

  Future<void> fetchProcurements() async {
    final response = await http.get(Uri.parse('http://localhost:8080/procurements'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        procurements = data.map((procurement) => Procurement.fromJson(procurement)).toList();
      });
    }
  }

  Future<void> createProcurement(Procurement procurement) async {
    final response = await http.post(
      Uri.parse('http://localhost:8080/procurements'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(procurement.toJson()),
    );
    if (response.statusCode == 200) {
      fetchProcurements();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Procurements'),
      ),
      body: ListView.builder(
        itemCount: procurements.length,
        itemBuilder: (context, index) {
          return ListTile(
            title: Text(procurements[index].title),
            subtitle: Text(procurements[index].description),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
            context: context,
            builder: (context) {
              return AlertDialog(
                title: Text('Create Procurement'),
                content: ProcurementForm(
                  onSubmit: (procurement) {
                    createProcurement(procurement);
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

class ProcurementForm extends StatefulWidget {
  final Function(Procurement) onSubmit;

  ProcurementForm({required this.onSubmit});

  @override
  _ProcurementFormState createState() => _ProcurementFormState();
}

class _ProcurementFormState extends State<ProcurementForm> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _statusController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextFormField(
            controller: _titleController,
            decoration: InputDecoration(labelText: 'Title'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a title';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _descriptionController,
            decoration: InputDecoration(labelText: 'Description'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a description';
              }
              return null;
            },
          ),
          TextFormField(
            controller: _statusController,
            decoration: InputDecoration(labelText: 'Status'),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a status';
              }
              return null;
            },
          ),
          ElevatedButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                widget.onSubmit(
                  Procurement(
                    title: _titleController.text,
                    description: _descriptionController.text,
                    status: _statusController.text,
                    dateCreated: DateTime.now(),
                    dateUpdated: DateTime.now(),
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