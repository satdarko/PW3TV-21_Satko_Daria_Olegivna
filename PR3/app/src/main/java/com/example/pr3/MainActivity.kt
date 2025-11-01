package com.example.pr3

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                SolarProfitCalculator()
            }
        }
    }
}

@Composable
fun SolarProfitCalculator() {
    var power by remember { mutableStateOf("10.0") }
    var hours by remember { mutableStateOf("5.5") }
    var tariff by remember { mutableStateOf("6.0") }
    var efficiency by remember { mutableStateOf("0.85") }
    var predictedPower by remember { mutableStateOf("52.0") }

    var result by remember { mutableStateOf<ResultData?>(null) }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF9FAFB))
            .padding(16.dp)
            .verticalScroll(rememberScrollState())
    ) {
        Text("Калькулятор прибутку СЕС", fontWeight = FontWeight.Bold, fontSize = 22.sp)
        Text("Введіть параметри станції та прогноз", color = Color.Gray, fontSize = 14.sp)

        Spacer(Modifier.height(16.dp))

        InputField("Потужність станції (кВт)", power) { power = it }
        InputField("Сонячні години (год/доба)", hours) { hours = it }
        InputField("Тариф (грн/кВт·год)", tariff) { tariff = it }
        InputField("Ефективність (0..1)", efficiency) { efficiency = it }
        InputField("Прогнозована потужність (кВт·год)", predictedPower) { predictedPower = it }

        Spacer(Modifier.height(20.dp))

        Button(
            onClick = {
                val res = calculateProfit(power, hours, tariff, efficiency, predictedPower)
                result = res
            },
            modifier = Modifier.fillMaxWidth().height(50.dp),
            shape = RoundedCornerShape(10.dp)
        ) {
            Text("Розрахувати", fontSize = 18.sp)
        }

        result?.let {
            Spacer(Modifier.height(20.dp))
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(10.dp),
                elevation = CardDefaults.cardElevation(5.dp)
            ) {
                Column(Modifier.padding(16.dp)) {
                    Text("Результати:", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Spacer(Modifier.height(8.dp))
                    Text("Добовий прибуток: %.2f грн".format(it.daily))
                    Text("Прогнозований прибуток: %.2f грн".format(it.predicted))
                    Divider(Modifier.padding(vertical = 8.dp))
                    Text("Сумарний прибуток: %.2f грн".format(it.total), fontWeight = FontWeight.Bold)
                }
            }
        }

        Spacer(Modifier.height(16.dp))
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(containerColor = Color(0xFFEFF6FF))
        ) {
            Column(Modifier.padding(12.dp)) {
                Text("Примітка:", fontWeight = FontWeight.Bold)
                Text(
                    "- Використовується прогноз сонячної потужності для наступного дня.\n" +
                            "- Тариф береться за поточним «зеленим» тарифом.\n" +
                            "- Усі значення вводяться у відповідних одиницях.",
                    fontSize = 13.sp
                )
            }
        }
    }
}

@Composable
fun InputField(label: String, value: String, onChange: (String) -> Unit) {
    OutlinedTextField(
        value = value,
        onValueChange = onChange,
        label = { Text(label) },
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        singleLine = true
    )
}

data class ResultData(val daily: Double, val predicted: Double, val total: Double)

fun calculateProfit(powerText: String, hoursText: String, tariffText: String, effText: String, predText: String): ResultData {
    val power = powerText.toDoubleOrNull() ?: 0.0
    val hours = hoursText.toDoubleOrNull() ?: 0.0
    val tariff = tariffText.toDoubleOrNull() ?: 0.0
    val eff = effText.toDoubleOrNull()?.coerceIn(0.0, 1.0) ?: 0.0
    val pred = predText.toDoubleOrNull() ?: 0.0

    val daily = power * hours * tariff * eff
    val predicted = pred * tariff
    val total = daily + predicted

    return ResultData(daily, predicted, total)
}
