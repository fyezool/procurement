import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:procurement/main.dart';
import 'package:procurement/screens/procurement_list_screen.dart';
import 'package:procurement/screens/supplier_list_screen.dart';
import 'package:procurement/screens/item_list_screen.dart';
import 'package:procurement/screens/quote_list_screen.dart';
import 'package:procurement/screens/contract_list_screen.dart';

void main() {
  testWidgets('HomeScreen navigation test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(MyApp());

    // Verify that the HomeScreen is displayed.
    expect(find.text('Procurement App'), findsOneWidget);

    // Tap on the Procurements button and trigger a frame.
    await tester.tap(find.text('Procurements'));
    await tester.pumpAndSettle();

    // Verify that the ProcurementListScreen is displayed.
    expect(find.byType(ProcurementListScreen), findsOneWidget);

    // Navigate back to the HomeScreen.
    await tester.tap(find.byIcon(Icons.arrow_back));
    await tester.pumpAndSettle();

    // Tap on the Suppliers button and trigger a frame.
    await tester.tap(find.text('Suppliers'));
    await tester.pumpAndSettle();

    // Verify that the SupplierListScreen is displayed.
    expect(find.byType(SupplierListScreen), findsOneWidget);

    // Navigate back to the HomeScreen.
    await tester.tap(find.byIcon(Icons.arrow_back));
    await tester.pumpAndSettle();

    // Tap on the Items button and trigger a frame.
    await tester.tap(find.text('Items'));
    await tester.pumpAndSettle();

    // Verify that the ItemListScreen is displayed.
    expect(find.byType(ItemListScreen), findsOneWidget);

    // Navigate back to the HomeScreen.
    await tester.tap(find.byIcon(Icons.arrow_back));
    await tester.pumpAndSettle();

    // Tap on the Quotes button and trigger a frame.
    await tester.tap(find.text('Quotes'));
    await tester.pumpAndSettle();

    // Verify that the QuoteListScreen is displayed.
    expect(find.byType(QuoteListScreen), findsOneWidget);

    // Navigate back to the HomeScreen.
    await tester.tap(find.byIcon(Icons.arrow_back));
    await tester.pumpAndSettle();

    // Tap on the Contracts button and trigger a frame.
    await tester.tap(find.text('Contracts'));
    await tester.pumpAndSettle();

    // Verify that the ContractListScreen is displayed.
    expect(find.byType(ContractListScreen), findsOneWidget);
  });

  testWidgets('ProcurementListScreen test', (WidgetTester tester) async {
    // Build the ProcurementListScreen and trigger a frame.
    await tester.pumpWidget(MaterialApp(home: ProcurementListScreen()));

    // Verify that the ProcurementListScreen is displayed.
    expect(find.text('Procurements'), findsOneWidget);

    // Tap on the floating action button to open the form.
    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();

    // Verify that the form is displayed.
    expect(find.text('Create Procurement'), findsOneWidget);
  });

  testWidgets('SupplierListScreen test', (WidgetTester tester) async {
    // Build the SupplierListScreen and trigger a frame.
    await tester.pumpWidget(MaterialApp(home: SupplierListScreen()));

    // Verify that the SupplierListScreen is displayed.
    expect(find.text('Suppliers'), findsOneWidget);

    // Tap on the floating action button to open the form.
    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();

    // Verify that the form is displayed.
    expect(find.text('Create Supplier'), findsOneWidget);
  });

  testWidgets('ItemListScreen test', (WidgetTester tester) async {
    // Build the ItemListScreen and trigger a frame.
    await tester.pumpWidget(MaterialApp(home: ItemListScreen()));

    // Verify that the ItemListScreen is displayed.
    expect(find.text('Items'), findsOneWidget);

    // Tap on the floating action button to open the form.
    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();

    // Verify that the form is displayed.
    expect(find.text('Create Item'), findsOneWidget);
  });

  testWidgets('QuoteListScreen test', (WidgetTester tester) async {
    // Build the QuoteListScreen and trigger a frame.
    await tester.pumpWidget(MaterialApp(home: QuoteListScreen()));

    // Verify that the QuoteListScreen is displayed.
    expect(find.text('Quotes'), findsOneWidget);

    // Tap on the floating action button to open the form.
    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();

    // Verify that the form is displayed.
    expect(find.text('Create Quote'), findsOneWidget);
  });

  testWidgets('ContractListScreen test', (WidgetTester tester) async {
    // Build the ContractListScreen and trigger a frame.
    await tester.pumpWidget(MaterialApp(home: ContractListScreen()));

    // Verify that the ContractListScreen is displayed.
    expect(find.text('Contracts'), findsOneWidget);

    // Tap on the floating action button to open the form.
    await tester.tap(find.byType(FloatingActionButton));
    await tester.pumpAndSettle();

    // Verify that the form is displayed.
    expect(find.text('Create Contract'), findsOneWidget);
  });
}