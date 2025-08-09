import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../services/auth_service.dart';

class RequisitionsScreen extends StatelessWidget {
  const RequisitionsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    // Access the AuthService to check the user's role
    final authService = Provider.of<AuthService>(context, listen: false);
    final userRole = authService.user?['role'];

    return Scaffold(
      appBar: AppBar(
        title: const Text('Requisitions Hub'),
      ),
      body: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                ElevatedButton.icon(
                  icon: const Icon(Icons.list_alt),
                  label: const Text('View My Requisitions'),
                  onPressed: () {
                    context.go('/procurement/requisitions/my');
                  },
                  style: ElevatedButton.styleFrom(
                    minimumSize: const Size(280, 50),
                    textStyle: const TextStyle(fontSize: 18),
                  ),
                ),
                const SizedBox(height: 20),
                // Conditionally show the "View All" button for Admins
                if (userRole == 'Admin')
                  ElevatedButton.icon(
                    icon: const Icon(Icons.grid_view_sharp),
                    label: const Text('View All Requisitions'),
                    onPressed: () {
                      context.go('/admin/requisitions');
                    },
                    style: ElevatedButton.styleFrom(
                      minimumSize: const Size(280, 50),
                      textStyle: const TextStyle(fontSize: 18),
                    ),
                  ),
                const SizedBox(height: 20),
                ElevatedButton.icon(
                  icon: const Icon(Icons.add_circle_outline),
                  label: const Text('Create New Requisition'),
                  onPressed: () {
            context.go('/procurement/requisitions/create');
                  },
                  style: ElevatedButton.styleFrom(
                    minimumSize: const Size(280, 50),
                    textStyle: const TextStyle(fontSize: 18),
                    backgroundColor: Theme.of(context).colorScheme.primary,
                    foregroundColor: Theme.of(context).colorScheme.onPrimary,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
