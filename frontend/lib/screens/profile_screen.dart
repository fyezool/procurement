import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/user.dart';
import '../services/api_service.dart';
import '../services/auth_service.dart';
import '../widgets/empty_state_widget.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({Key? key}) : super(key: key);

  @override
  _ProfileScreenState createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  final _apiService = ApiService();
  final _nameFormKey = GlobalKey<FormState>();
  final _passwordFormKey = GlobalKey<FormState>();

  late TextEditingController _nameController;
  final _oldPasswordController = TextEditingController();
  final _newPasswordController = TextEditingController();

  late Future<User> _userFuture;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController();
    _userFuture = _apiService.getMyProfile();
    _userFuture.then((user) {
      if (mounted) {
        _nameController.text = user.name;
      }
    });
  }

  @override
  void dispose() {
    _nameController.dispose();
    _oldPasswordController.dispose();
    _newPasswordController.dispose();
    super.dispose();
  }

  void _updateProfile() async {
    if (_nameFormKey.currentState!.validate()) {
      try {
        await _apiService.updateMyProfile(_nameController.text);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Profile updated successfully'), backgroundColor: Colors.green),
        );
        // Refresh the user in AuthService to update the name in the AppBar
        Provider.of<AuthService>(context, listen: false).updateUserName(_nameController.text);
        setState(() {
          _userFuture = _apiService.getMyProfile();
        });
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to update profile: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  void _changePassword() async {
    if (_passwordFormKey.currentState!.validate()) {
       try {
        await _apiService.changeMyPassword(_oldPasswordController.text, _newPasswordController.text);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Password changed successfully'), backgroundColor: Colors.green),
        );
        _oldPasswordController.clear();
        _newPasswordController.clear();
        // Force re-login for security
        Provider.of<AuthService>(context, listen: false).logout();

      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to change password: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Profile'),
      ),
      body: FutureBuilder<User>(
        future: _userFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return EmptyStateWidget(
              message: 'Failed to load your profile: ${snapshot.error}',
              icon: Icons.error_outline,
              onRetry: () {
                setState(() {
                  _userFuture = _apiService.getMyProfile();
                });
              },
            );
          }
          final user = snapshot.data!;
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Email: ${user.email}', style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: 8),
                Text('Role: ${user.role}', style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: 24),
                const Divider(),
                const SizedBox(height: 16),
                Text('Update Profile', style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 16),
                Form(
                  key: _nameFormKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _nameController,
                        decoration: const InputDecoration(labelText: 'Name', border: OutlineInputBorder()),
                        validator: (value) => value!.isEmpty ? 'Please enter a name' : null,
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _updateProfile,
                        child: const Text('Update Name'),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 24),
                const Divider(),
                const SizedBox(height: 16),
                Text('Change Password', style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 16),
                 Form(
                  key: _passwordFormKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _oldPasswordController,
                        decoration: const InputDecoration(labelText: 'Old Password', border: OutlineInputBorder()),
                        obscureText: true,
                         validator: (value) => value!.isEmpty ? 'Please enter your old password' : null,
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _newPasswordController,
                        decoration: const InputDecoration(labelText: 'New Password', border: OutlineInputBorder()),
                        obscureText: true,
                        validator: (value) => value!.length < 8 ? 'Password must be at least 8 characters' : null,
                      ),
                       const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _changePassword,
                        child: const Text('Change Password'),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
