import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../models/vendor.dart';
import '../../services/api_service.dart';

class CreateRequisitionScreen extends StatefulWidget {
  const CreateRequisitionScreen({Key? key}) : super(key: key);

  @override
  _CreateRequisitionScreenState createState() =>
      _CreateRequisitionScreenState();
}

class _CreateRequisitionScreenState extends State<CreateRequisitionScreen> {
  final _formKey = GlobalKey<FormState>();
  final _apiService = ApiService();
  bool _isLoading = false;

  int? _selectedVendorId;
  final _itemDescriptionController = TextEditingController();
  final _quantityController = TextEditingController();
  final _estimatedPriceController = TextEditingController();
  final _justificationController = TextEditingController();

  late Future<List<Vendor>> _vendorsFuture;

  @override
  void initState() {
    super.initState();
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

  void _submit() async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });

      try {
        final requisitionData = {
          'vendor_id': _selectedVendorId,
          'item_description': _itemDescriptionController.text,
          'quantity': int.tryParse(_quantityController.text) ?? 0,
          'estimated_price': double.tryParse(_estimatedPriceController.text) ?? 0.0,
          'justification': _justificationController.text,
        };
        await _apiService.createRequisition(requisitionData);

        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
                content: Text('Requisition created successfully'),
                backgroundColor: Colors.green),
          );
          context.pop();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
                content: Text('Failed to create requisition: $e'),
                backgroundColor: Colors.red),
          );
        }
      } finally {
        if (mounted) {
          setState(() {
            _isLoading = false;
          });
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create New Requisition'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: SingleChildScrollView(
            child: FutureBuilder<List<Vendor>>(
              future: _vendorsFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }
                if (snapshot.hasError || !snapshot.hasData || snapshot.data!.isEmpty) {
                  return const Center(
                      child: Text('Could not load vendors. Please add a vendor first.'));
                }

                final vendors = snapshot.data!;
                return Column(
                  children: [
                    DropdownButtonFormField<int>(
                      value: _selectedVendorId,
                      decoration: const InputDecoration(
                          labelText: 'Vendor', border: OutlineInputBorder()),
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
                      validator: (value) =>
                          value == null ? 'Please select a vendor' : null,
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _itemDescriptionController,
                      decoration: const InputDecoration(
                          labelText: 'Item Description', border: OutlineInputBorder()),
                      validator: (value) =>
                          value!.isEmpty ? 'Please enter a description' : null,
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _quantityController,
                      decoration: const InputDecoration(
                          labelText: 'Quantity', border: OutlineInputBorder()),
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
                      decoration: const InputDecoration(
                          labelText: 'Estimated Price (\$)', border: OutlineInputBorder()),
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
                      decoration: const InputDecoration(
                          labelText: 'Justification', border: OutlineInputBorder()),
                      maxLines: 3,
                    ),
                    const SizedBox(height: 20),
                    _isLoading
                        ? const CircularProgressIndicator()
                        : ElevatedButton(
                            style: ElevatedButton.styleFrom(
                              minimumSize: const Size(double.infinity, 50),
                            ),
                            onPressed: _submit,
                            child: const Text('Submit Requisition'),
                          ),
                  ],
                );
              },
            ),
          ),
        ),
      ),
    );
  }
}
