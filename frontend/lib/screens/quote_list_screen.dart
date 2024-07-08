import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../proc/quote.dart';

class QuoteListScreen extends StatefulWidget {
  @override
  _QuoteListScreenState createState() => _QuoteListScreenState();
}

class _QuoteListScreenState extends State<QuoteListScreen> {
  List<Quote> quotes = [];
  List<String> procurementIds = [];
  List<String> supplierIds = [];

  @override
  void initState() {
    super.initState();
    fetchQuotes();
    fetchProcurementIds();
    fetchSupplierIds();
  }

  Future<void> fetchQuotes() async {
    final response = await http.get(Uri.parse('http://localhost:8080/quotes'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        quotes = data.map((quote) => Quote.fromJson(quote)).toList();
      });
    }
  }

  Future<void> fetchProcurementIds() async {
    final response =
        await http.get(Uri.parse('http://localhost:8080/procurements'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        procurementIds = data.map((item) => item['id'].toString()).toList();
      });
    }
  }

  Future<void> fetchSupplierIds() async {
    final response =
        await http.get(Uri.parse('http://localhost:8080/suppliers'));
    if (response.statusCode == 200) {
      List<dynamic> data = jsonDecode(response.body);
      setState(() {
        supplierIds = data.map((item) => item['id'].toString()).toList();
      });
    }
  }

  Future<void> createQuote(Quote quote) async {
    final response = await http.post(
      Uri.parse('http://localhost:8080/quotes'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(quote.toJson()),
    );
    if (response.statusCode == 200) {
      fetchQuotes();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Quotes'),
      ),
      body: ListView.builder(
        itemCount: quotes.length,
        itemBuilder: (context, index) {
          return ListTile(
            title: Text('Quote ${quotes[index].id}'),
            subtitle: Text('Total Cost: \$${quotes[index].totalCost}'),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
            context: context,
            builder: (context) {
              return AlertDialog(
                title: Text('Create Quote'),
                content: QuoteForm(
                  onSubmit: (quote) {
                    createQuote(quote);
                    Navigator.of(context).pop();
                  },
                  procurementIds: procurementIds,
                  supplierIds: supplierIds,
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

class QuoteForm extends StatefulWidget {
  final Function(Quote) onSubmit;
  final List<String> procurementIds;
  final List<String> supplierIds;

  QuoteForm({
    required this.onSubmit,
    required this.procurementIds,
    required this.supplierIds,
  });

  @override
  _QuoteFormState createState() => _QuoteFormState();
}

class _QuoteFormState extends State<QuoteForm> {
  final _formKey = GlobalKey<FormState>();
  String? _selectedProcurementId;
  String? _selectedSupplierId;
  final _totalCostController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          DropdownButtonFormField<String>(
            value: _selectedProcurementId,
            decoration: InputDecoration(labelText: 'Procurement ID'),
            items: widget.procurementIds.map((String id) {
              return DropdownMenuItem<String>(
                value: id,
                child: Text(id),
              );
            }).toList(),
            onChanged: (String? newValue) {
              setState(() {
                _selectedProcurementId = newValue;
              });
            },
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please select a procurement ID';
              }
              return null;
            },
          ),
          DropdownButtonFormField<String>(
            value: _selectedSupplierId,
            decoration: InputDecoration(labelText: 'Supplier ID'),
            items: widget.supplierIds.map((String id) {
              return DropdownMenuItem<String>(
                value: id,
                child: Text(id),
              );
            }).toList(),
            onChanged: (String? newValue) {
              setState(() {
                _selectedSupplierId = newValue;
              });
            },
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please select a supplier ID';
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
                  Quote(
                    procurementId: _selectedProcurementId!,
                    supplierId: _selectedSupplierId!,
                    items: [],
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
