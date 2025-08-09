import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'screens/login_screen.dart';
import 'screens/dashboard_screen.dart';
import 'screens/vendors_screen.dart';
import 'screens/purchase_orders_screen.dart';
import 'services/auth_service.dart';
import 'services/api_service.dart';

void main() {
  runApp(const ProcurementApp());
}

class ProcurementApp extends StatelessWidget {
  const ProcurementApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthService()),
        Provider(create: (_) => ApiService()),
      ],
      child: Consumer<AuthService>(
        builder: (context, authService, child) {
          final router = GoRouter(
            initialLocation: authService.isAuthenticated ? '/dashboard' : '/login',
            routes: [
              GoRoute(
                path: '/login',
                builder: (context, state) => const LoginScreen(),
              ),
              GoRoute(
                path: '/dashboard',
                builder: (context, state) => const DashboardScreen(),
              ),
              GoRoute(
                path: '/vendors',
                builder: (context, state) => const VendorsScreen(),
              ),
              GoRoute(
                path: '/purchase-orders',
                builder: (context, state) => const PurchaseOrdersScreen(),
              ),
            ],
            redirect: (context, state) {
              final isAuthenticated = authService.isAuthenticated;
              final isLoggingIn = state.uri.path == '/login';

              if (!isAuthenticated && !isLoggingIn) return '/login';
              if (isAuthenticated && isLoggingIn) return '/dashboard';
              return null;
            },
          );

          return MaterialApp.router(
            title: 'Procurement Management System',
            theme: ThemeData(
              colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
              useMaterial3: true,
            ),
            routerConfig: router,
            debugShowCheckedModeBanner: false,
          );
        },
      ),
    );
  }
}
