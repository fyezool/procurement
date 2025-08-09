import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../services/auth_service.dart';
import '../services/api_service.dart';
import '../models/navigation_item.dart';
import '../widgets/breadcrumbs.dart';


class AppShell extends StatefulWidget {
  final Widget child;

  const AppShell({super.key, required this.child});

  @override
  _AppShellState createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  List<NavigationItem> _menuItems = [];
  List<NavigationItem> _filteredMenuItems = [];
  bool _isLoading = true;
  final TextEditingController _searchController = TextEditingController();
  bool _isLoading = true;


  @override
  void initState() {
    super.initState();
    _fetchMenu();
    _searchController.addListener(_filterMenu);
  }

  @override
  void dispose() {
    _searchController.removeListener(_filterMenu);
    _searchController.dispose();
    super.dispose();

  }

  Future<void> _fetchMenu() async {
    try {
      final response = await ApiService.get('/navigation/menu');
      final List<dynamic> menuData = response.data;
      setState(() {
        _menuItems = menuData.map((item) => NavigationItem.fromJson(item)).toList();
        _filteredMenuItems = _menuItems;

        _isLoading = false;
      });
    } catch (e) {
      print('Failed to fetch menu: $e');
      setState(() {
        _isLoading = false;
      });
    }
  }

  void _filterMenu() {
    final query = _searchController.text.toLowerCase();
    if (query.isEmpty) {
      setState(() {
        _filteredMenuItems = _menuItems;
      });
      return;
    }

    final filtered = _menuItems.where((item) {
      final titleMatch = item.title.toLowerCase().contains(query);
      final subItemMatch = item.subItems.any((sub) => sub.title.toLowerCase().contains(query));
      return titleMatch || subItemMatch;
    }).toList();

    setState(() {
      _filteredMenuItems = filtered;
    });
  }

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
          IconButton(
            icon: Badge(
              label: Text('3'), // Static badge for now
              child: Icon(Icons.notifications),
            ),
            onPressed: () {
              // TODO: Implement notification panel
            },
          ),
          const SizedBox(width: 8),
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

              }
            },
          ),
          const SizedBox(width: 16),
        ],
      ),

      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // TODO: Implement quick actions
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Quick Actions coming soon!')),
          );
        },
        child: const Icon(Icons.add),
      ),
      body: Row(
        children: [
          if (_isLoading)
            const Center(child: CircularProgressIndicator())
          else
            NavigationDrawer(
              children: [
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: TextField(
                    controller: _searchController,
                    decoration: const InputDecoration(
                      hintText: 'Search menu...',
                      prefixIcon: Icon(Icons.search),
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                ..._filteredMenuItems.map((item) {
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
                }),
              ],
            ),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Breadcrumbs(),
                Expanded(child: widget.child),
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
      case 'description':
        return Icons.description;
      case 'check_circle':
        return Icons.check_circle;
      case 'list_alt':
        return Icons.list_alt;
      case 'receipt':
        return Icons.receipt;
      default:
        return Icons.circle;
    }
  }
}
