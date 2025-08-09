import 'package:flutter/material.dart';
import '../models/requisition.dart';
import '../models/vendor.dart';
import '../services/api_service.dart';

class EditRequisitionDialog extends StatefulWidget {
  final Requisition requisition;
  final Function(Map<String, dynamic>) onSave;

  const EditRequisitionDialog({
    Key? key,
    required this.requisition,
    required this.onSave,
  }) : super(key: key);

  @override
  _EditRequisitionDialogState createState() => _EditRequisitionDialogState();
}

class _EditRequisitionDialogState extends State<EditRequisitionDialog> {
  final _formKey = GlobalKey<FormState>();
  final _apiService = ApiService();

  late TextEditingController _itemDescriptionController;
  late TextEditingController _quantityController;
  late TextEditingController _estimatedPriceController;
  late TextEditingController _justificationController;
  int? _selectedVendorId;

  late Future<List<Vendor>> _vendorsFuture;

  @override
  void initState() {
    super.initState();
    _itemDescriptionController = TextEditingController(text: widget.requisition.itemDescription);
    _quantityController = TextEditingController(text: widget.requisition.quantity.toString());
    _estimatedPriceController = TextEditingController(text: widget.requisition.estimatedPrice.toString());
    _justificationController = TextEditingController(text: widget.requisition.justification);
    _selectedVendorId = widget.requisition.vendorId;
    _vendorsFuture = _apiService.getVendors();
  }

  @override
  void dispose() {
    _itemDescriptionController.dispose();
    _quantityController.dispose();
    _estimatedPriceController.dispose();
    _justificationController.dispose();
    super.dispose();
  }

  void _submit() {
    if (_formKey.currentState!.validate()) {
      final updatedData = {
        'vendor_id': _selectedVendorId,
        'item_description': _itemDescriptionController.text,
        'quantity': int.tryParse(_quantityController.text) ?? 0,
        'estimated_price': double.tryParse(_estimatedPriceController.text) ?? 0.0,
        'justification': _justificationController.text,
      };
      widget.onSave(updatedData);
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Edit Requisition'),
      content: FutureBuilder<List<Vendor>>(
        future: _vendorsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
           if (snapshot.hasError) {
            return Center(child: Text("Error loading vendors: ${snapshot.error}"));
          }
          if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text("No vendors available. Please add a vendor first."));
          }

          final vendors = snapshot.data!;
          return Form(
            key: _formKey,
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  DropdownButtonFormField<int>(
                    value: _selectedVendorId,
                    decoration: const InputDecoration(labelText: 'Vendor', border: OutlineInputBorder()),
                    items: vendors.map((Vendor vendor) {
                      return DropdownMenuItem<int>(
                        value: vendor.id,
                        child: Text(vendor.name),
                      );
                    }).toList(),
                    onChanged: (newValue) {
                      setState(() {
                        _selectedVendorId = newValue;
                      });
                    },
                    validator: (value) => value == null ? 'Please select a vendor' : null,
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _itemDescriptionController,
                    decoration: const InputDecoration(labelText: 'Item Description', border: OutlineInputBorder()),
                    validator: (value) => value!.isEmpty ? 'Please enter a description' : null,
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _quantityController,
                    decoration: const InputDecoration(labelText: 'Quantity', border: OutlineInputBorder()),
                    keyboardType: TextInputType.number,
                    validator: (value) {
                       if (value == null || value.isEmpty) return 'Please enter a quantity';
                       if (int.tryParse(value) == null) return 'Please enter a valid number';
                       return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _estimatedPriceController,
                    decoration: const InputDecoration(labelText: 'Estimated Price', border: OutlineInputBorder()),
                    keyboardType: const TextInputType.numberWithOptions(decimal: true),
                     validator: (value) {
                       if (value == null || value.isEmpty) return 'Please enter a price';
                       if (double.tryParse(value) == null) return 'Please enter a valid price';
                       return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _justificationController,
                    decoration: const InputDecoration(labelText: 'Justification', border: OutlineInputBorder()),
                  ),
                ],
              ),
            ),
          );
        },
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _submit,
          child: const Text('Save'),
        ),
      ],
    );
  }
}
