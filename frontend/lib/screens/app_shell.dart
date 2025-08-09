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
  int _bottomNavIndex = 0;

  @override
  void initState() {
    super.initState();
    _fetchMenu();
    _searchController.addListener(_filterMenu);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _updateBottomNavIndex();
  }

  void _updateBottomNavIndex() {
    final path = GoRouter.of(context).routerDelegate.currentConfiguration.uri.toString();
    final topLevelItems = _menuItems.where((item) => item.subItems.isEmpty).toList();
    final index = topLevelItems.indexWhere((item) => item.path == path);
    if (index != -1 && _bottomNavIndex != index) {
      setState(() {
        _bottomNavIndex = index;
      });
    }
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

  Widget _buildSideDrawer() {
    return NavigationDrawer(
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
    );
  }

  @override
  Widget build(BuildContext context) {
    final authService = Provider.of<AuthService>(context);

    return LayoutBuilder(
      builder: (context, constraints) {
        final isMobile = constraints.maxWidth < 600;

        return Scaffold(
          appBar: AppBar(
            title: const Text('Procurement System'),
            actions: [
              IconButton(
                icon: Badge(
                  label: Text('3'), // Static badge for now
                  child: Icon(Icons.notifications),
                ),
                onPressed: () {},
              ),
              const SizedBox(width: 8),
              PopupMenuButton(
                icon: CircleAvatar(
                  child: Text(authService.user?['name']?.substring(0, 1) ?? 'U'),
                ),
                itemBuilder: (context) => [
                  const PopupMenuItem(value: 'profile', child: Text('Profile')),
                  const PopupMenuItem(value: 'logout', child: Text('Logout')),
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
          drawer: isMobile ? _buildSideDrawer() : null,
          body: Row(
            children: [
              if (!isMobile) _buildSideDrawer(),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Breadcrumbs(),
                    Expanded(child: widget.child),
                  ],
                ),
              ),
            ],
          ),
          bottomNavigationBar: isMobile ? _buildBottomNav() : null,
          floatingActionButton: FloatingActionButton(
            onPressed: () {},
            child: const Icon(Icons.add),
          ),
        );
      },
    );
  }

  Widget _buildBottomNav() {
    final topLevelItems = _menuItems.where((item) => item.subItems.isEmpty).toList();

    return BottomNavigationBar(
      currentIndex: _bottomNavIndex,
      onTap: (index) {
        setState(() {
          _bottomNavIndex = index;
        });
        context.go(topLevelItems[index].path);
      },
      items: topLevelItems.map((item) {
        return BottomNavigationBarItem(
          icon: Icon(getIconData(item.icon)),
          label: item.title,
        );
      }).toList(),
    );
  }

  IconData getIconData(String iconName) {
    switch (iconName) {
      case 'dashboard': return Icons.dashboard;
      case 'shopping_cart': return Icons.shopping_cart;
      case 'store': return Icons.store;
      case 'assessment': return Icons.assessment;
      case 'settings': return Icons.settings;
      case 'description': return Icons.description;
      case 'check_circle': return Icons.check_circle;
      case 'list_alt': return Icons.list_alt;
      case 'receipt': return Icons.receipt;
      default: return Icons.circle;
    }
  }
}
