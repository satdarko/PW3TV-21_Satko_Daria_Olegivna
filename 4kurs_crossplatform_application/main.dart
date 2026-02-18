import 'package:flutter/material.dart';
import 'dart:math';

void main() {
  runApp(const SolarApp());
}

class SolarApp extends StatelessWidget {
  const SolarApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Solar Calculator',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const CalculationScreen(),
    );
  }
}

class CalculationScreen extends StatefulWidget {
  const CalculationScreen({super.key});

  @override
  State<CalculationScreen> createState() => _CalculationScreenState();
}

class _CalculationScreenState extends State<CalculationScreen> {
  final powerCtrl = TextEditingController(text: "5.0");
  final sigma1Ctrl = TextEditingController(text: "1.0");
  final sigma2Ctrl = TextEditingController(text: "0.25");
  final costCtrl = TextEditingController(text: "7.0");

  String result = "";

  double erf(double x) {
    const double a1 = 0.254829592;
    const double a2 = -0.284496736;
    const double a3 = 1.421413741;
    const double a4 = -1.453152027;
    const double a5 = 1.061405429;
    const double p = 0.3275911;

    int sign = 1;
    if (x < 0) sign = -1;
    x = x.abs();

    double t = 1.0 / (1.0 + p * x);
    double y = 1.0 - (((((a5 * t + a4) * t) + a3) * t + a2) * t + a1) * t * exp(-x * x);

    return sign * y;
  }

  double calculateProfit(double pc, double sigma, double price) {
    double delta = pc * 0.05;
    double p1 = pc - delta;
    double p2 = pc + delta;

    double integral = 0.5 * (erf((p2 - pc) / (sigma * sqrt(2))) - erf((p1 - pc) / (sigma * sqrt(2))));

    double w1 = pc * 24 * integral;
    double w2 = pc * 24 * (1 - integral);

    double profit = (w1 * price * 1000) - (w2 * price * 1000);
    return profit / 1000;
  }

  void calculate() {
    double pc = double.tryParse(powerCtrl.text) ?? 0;
    double s1 = double.tryParse(sigma1Ctrl.text) ?? 1;
    double s2 = double.tryParse(sigma2Ctrl.text) ?? 1;
    double b = double.tryParse(costCtrl.text) ?? 0;

    double profit1 = calculateProfit(pc, s1, b);
    double profit2 = calculateProfit(pc, s2, b);
    double diff = profit2 - profit1;

    setState(() {
      result = "Прибуток (система 1): ${profit1.toStringAsFixed(2)} тис. грн\n"
          "Прибуток (система 2): ${profit2.toStringAsFixed(2)} тис. грн\n"
          "Виграш: ${diff.toStringAsFixed(2)} тис. грн";
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Розрахунок прибутку СЕС")),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextField(
              controller: powerCtrl,
              decoration: const InputDecoration(labelText: "Потужність (Pc), МВт", border: OutlineInputBorder()),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 10),
            TextField(
              controller: sigma1Ctrl,
              decoration: const InputDecoration(labelText: "Похибка 1 (sigma1)", border: OutlineInputBorder()),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 10),
            TextField(
              controller: sigma2Ctrl,
              decoration: const InputDecoration(labelText: "Похибка 2 (sigma2)", border: OutlineInputBorder()),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 10),
            TextField(
              controller: costCtrl,
              decoration: const InputDecoration(labelText: "Вартість (B), грн/кВт·год", border: OutlineInputBorder()),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: calculate,
              child: const Text("Розрахувати"),
            ),
            const SizedBox(height: 20),
            if (result.isNotEmpty)
              Container(
                padding: const EdgeInsets.all(16),
                color: Colors.blue[50],
                child: Text(result, style: const TextStyle(fontSize: 16)),
              ),
          ],
        ),
      ),
    );
  }
}