import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../models/breadcrumb_item.dart';
import '../services/api_service.dart';

class Breadcrumbs extends StatefulWidget {
  const Breadcrumbs({Key? key}) : super(key: key);

  @override
  _BreadcrumbsState createState() => _BreadcrumbsState();
}

class _BreadcrumbsState extends State<Breadcrumbs> {
  List<BreadcrumbItem> _breadcrumbs = [];

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final path = GoRouter.of(context).routerDelegate.currentConfiguration.uri.toString();
    _fetchBreadcrumbs(path);
  }

  Future<void> _fetchBreadcrumbs(String path) async {
    if (path.isEmpty) return;
    try {
      final response = await ApiService.get('/navigation/breadcrumbs', queryParameters: {'path': path});
      final List<dynamic> data = response.data;
      setState(() {
        _breadcrumbs = data.map((item) => BreadcrumbItem.fromJson(item)).toList();
      });
    } catch (e) {
      print('Failed to fetch breadcrumbs: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_breadcrumbs.isEmpty) {
      return const SizedBox.shrink();
    }

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      child: Wrap(
        spacing: 8.0,
        runSpacing: 4.0,
        crossAxisAlignment: WrapCrossAlignment.center,
        children: _buildBreadcrumbWidgets(context),
      ),
    );
  }

  List<Widget> _buildBreadcrumbWidgets(BuildContext context) {
    final List<Widget> widgets = [];
    for (int i = 0; i < _breadcrumbs.length; i++) {
      final item = _breadcrumbs[i];
      final isLast = i == _breadcrumbs.length - 1;

      widgets.add(
        InkWell(
          onTap: isLast ? null : () => context.go(item.path),
          child: Text(
            item.title,
            style: TextStyle(
              color: isLast ? Colors.black : Colors.blue,
              fontWeight: isLast ? FontWeight.bold : FontWeight.normal,
            ),
          ),
        ),
      );

      if (!isLast) {
        widgets.add(const Icon(Icons.chevron_right, size: 16.0));
      }
    }
    return widgets;
  }
}
