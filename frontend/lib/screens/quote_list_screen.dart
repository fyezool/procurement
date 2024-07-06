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

  @override
  void initState() {
    super.initState();
    fetchQuotes();
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

  QuoteForm({required this.onSubmit});

  @override
  _QuoteFormState createState() => _QuoteFormState();
}

class _QuoteFormState extends State<QuoteForm> {
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
                  Quote(
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