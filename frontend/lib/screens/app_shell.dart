import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../services/auth_service.dart';
import '../services/api_service.dart';
import '../models/navigation_item.dart';

class AppShell extends StatefulWidget {
  final Widget child;

  const AppShell({super.key, required this.child});

  @override
  _AppShellState createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  List<NavigationItem> _menuItems = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchMenu();
  }

  Future<void> _fetchMenu() async {
    try {
      final response = await ApiService.get('/navigation/menu');
      final List<dynamic> menuData = response.data;
      setState(() {
        _menuItems = menuData.map((item) => NavigationItem.fromJson(item)).toList();
        _isLoading = false;
      });
    } catch (e) {
      print('Failed to fetch menu: $e');
      setState(() {
        _isLoading = false;
      });
      // Handle error, maybe show a snackbar
    }
  }

  @override
  Widget build(BuildContext context) {
    final authService = Provider.of<AuthService>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Procurement System'),
        actions: [
          PopupMenuButton(
            icon: CircleAvatar(
              child: Text(authService.user?['name']?.substring(0, 1) ?? 'U'),
            ),
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'profile',
                child: Text('Profile'),
              ),
              const PopupMenuItem(
                value: 'logout',
                child: Text('Logout'),
              ),
            ],
            onSelected: (value) {
              if (value == 'logout') {
                authService.logout();
                // The router's redirect will handle navigation to the login screen
              }
            },
          ),
          const SizedBox(width: 16),
        ],
      ),
      body: Row(
        children: [
          if (_isLoading)
            const CircularProgressIndicator()
          else
            NavigationDrawer(
              children: _menuItems.map((item) {
                if (item.subItems.isEmpty) {
                  return ListTile(
                    leading: Icon(getIconData(item.icon)),
                    title: Text(item.title),
                    onTap: () => context.go(item.path),
                  );
                } else {
                  return ExpansionTile(
                    leading: Icon(getIconData(item.icon)),
                    title: Text(item.title),
                    children: item.subItems.map((subItem) {
                      return ListTile(
                        title: Text(subItem.title),
                        onTap: () => context.go(subItem.path),
                      );
                    }).toList(),
                  );
                }
              }).toList(),
            ),
          Expanded(
            child: widget.child,
          ),
        ],
      ),
    );
  }

  IconData getIconData(String iconName) {
    switch (iconName) {
      case 'dashboard':
        return Icons.dashboard;
      case 'shopping_cart':
        return Icons.shopping_cart;
      case 'store':
        return Icons.store;
      case 'assessment':
        return Icons.assessment;
      case 'settings':
        return Icons.settings;
      default:
        return Icons.circle;
    }
  }
}
