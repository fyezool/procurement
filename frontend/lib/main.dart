import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';

// Screens
import 'screens/login_screen.dart';
import 'screens/dashboard_screen.dart';
import 'screens/app_shell.dart';
import 'screens/procurement/procurement_main_screen.dart';
import 'screens/procurement/requisitions_screen.dart';
import 'screens/procurement/create_requisition_screen.dart';
import 'screens/procurement/my_requisitions_screen.dart';
import 'screens/procurement/pending_requisitions_screen.dart';
import 'screens/purchase_orders_screen.dart';
import 'screens/procurement/create_purchase_order_screen.dart';
import 'screens/procurement/approvals_screen.dart';
import 'screens/vendors_screen.dart';
import 'screens/vendors/add_vendor_screen.dart';
import 'screens/vendors/vendor_performance_screen.dart';
import 'screens/vendors/vendor_communications_screen.dart';
import 'screens/reports/reports_screen.dart';
import 'screens/reports/analytics_screen.dart';
import 'screens/reports/vendor_reports_screen.dart';
import 'screens/reports/spend_analysis_screen.dart';
import 'screens/admin/admin_screen.dart';
import 'screens/admin/user_management_screen.dart';
import 'screens/admin/role_management_screen.dart';
import 'screens/admin/system_settings_screen.dart';
import 'screens/admin/all_requisitions_screen.dart';
import 'screens/admin/all_purchase_orders_screen.dart';

// Services
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
            refreshListenable: authService,
            initialLocation: '/login',
            routes: [
              GoRoute(
                path: '/login',
                builder: (context, state) => const LoginScreen(),
              ),
              ShellRoute(
                builder: (context, state, child) {
                  return AppShell(child: child);
                },
                routes: [
                  GoRoute(
                    path: '/dashboard',
                    builder: (context, state) => const DashboardScreen(),
                  ),
                  // Procurement Routes
                  GoRoute(
                    path: '/procurement',
                    builder: (context, state) => const ProcurementMainScreen(),
                  ),
                  GoRoute(
                    path: '/procurement/requisitions',
                    builder: (context, state) => const RequisitionsScreen(),
                  ),
                  GoRoute(
                    path: '/procurement/requisitions/my',
                    builder: (context, state) => const MyRequisitionsScreen(),
                  ),
                  GoRoute(
                    path: '/procurement/purchase-orders',
                    builder: (context, state) => const PurchaseOrdersScreen(),
                  ),
                   GoRoute(
                    path: '/procurement/approvals',
                    builder: (context, state) => const ApprovalsScreen(),
                  ),
                  // Vendor Routes
                  GoRoute(
                    path: '/vendors',
                    builder: (context, state) => const VendorsScreen(),
                  ),
                  // Reports Routes
                  GoRoute(
                    path: '/reports',
                    builder: (context, state) => const ReportsScreen(),
                  ),
                  // Admin Routes
                   GoRoute(
                    path: '/admin',
                    builder: (context, state) => const AdminScreen(),
                  ),
                   GoRoute(
                    path: '/admin/users',
                    builder: (context, state) => const UserManagementScreen(),
                  ),
                  GoRoute(
                    path: '/admin/requisitions',
                    builder: (context, state) => const AllRequisitionsScreen(),
                  ),
                  GoRoute(
                    path: '/admin/purchase-orders',
                    builder: (context, state) => const AllPurchaseOrdersScreen(),
                  ),
                   GoRoute(
                    path: '/admin/settings',
                    builder: (context, state) => const SystemSettingsScreen(),
                  ),
                ],
              ),
            ],
            redirect: (context, state) {
              final isAuthenticated = authService.isAuthenticated;
              final isLoggingIn = state.uri.path == '/login';

              if (!isAuthenticated && !isLoggingIn) {
                return '/login';
              }
              if (isAuthenticated && isLoggingIn) {
                return '/dashboard';
              }
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
