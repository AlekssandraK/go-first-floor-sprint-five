package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений
const (
	MInKm      = 1000 //количество метров в одном километре
	MinInHours = 60   //количество минут в одном часе
	LenStep    = 0.65 //длина одного шага
	CmInM      = 100  //количество сантиметров в одном метре
)

// Training - общая структура для всех тренировок
type Training struct {
	TrainingType string        //тип тренировки
	Action       int           //количество повторов(шаги, гребки при плавании)
	LenStep      float64       //длина одного шага или гребка в м
	Duration     time.Duration //продолжительность тренировки
	Weight       float64       //вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь
// Формула расчета: количество_поворотов * длина_шага м_в_км
func (t Training) distance() float64 {
	d := float64(t.Action) * t.LenStep / MInKm
	return d
}

// константа для сравнения числа типа float64 c 0
const eps = 1e-8

// meanSpeed возвращает среднюю скорость бега или ходьбы
func (t Training) meanSpeed() float64 {
	//проверка делителя на нулевое значение
	if math.Abs(t.Duration.Hours()) < eps {
		return 0
	}
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке
// Пока возвращаем 0, так как этот метод будет переопределяться для каждого типа тренрировки
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке
type InfoMessage struct {
	TrainingType string        //тип тренировки
	Duration     time.Duration //длительность тренировки
	Distance     float64       //расстояние, которое преодолел пользователь
	Speed        float64       //средняя скорость, с которой двигался пользователь
	Calories     float64       //количество потраченных килокалорий на тренировке
}

// TrainingInfo возвращает структуру InfoMessage, в которой хранится вся информация о проведенной тренировке
func (t Training) TrainingInfo() InfoMessage {
	im := InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
	return im
}

// String возвращает строку с информацией о проведенной тренировке
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км. \nСр.скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaoloriesCalculator интерфейс для структур: Running, Walking и Swimming
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потреченных килокалорий при беге
const (
	CaloriesMeanSpeedMultiplier = 18   //множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 //коэффициент изменения средней скорости
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training Training
}

// Calories возвращает количество потраченных килокалорий при беге
// Формула расчета: ((18 * средняя_скорость_в_км/ч + 1.79) * вес_спортсмена_в_кг/м_в_км * время_тренировки_в_часах * мин_в_часе)
// Это переопределенный метод Calories() из Training
func (r Running) Calories() float64 {
	calories := (CaloriesMeanSpeedMultiplier*r.Training.meanSpeed() +
		CaloriesMeanSpeedShift) * r.Training.Weight / MInKm *
		r.Training.Duration.Hours() * MinInHours
	return calories
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке
// Это переопределенный метод TrainingInfo() из Training
func (r Running) TrainingInfo() InfoMessage {
	im := InfoMessage{
		TrainingType: r.Training.TrainingType,
		Duration:     r.Training.Duration,
		Distance:     r.Training.distance(),
		Speed:        r.Training.meanSpeed(),
		Calories:     r.Training.Calories(),
	}
	return im
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 //коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 //коэффициент для роста
	KmHInMsec                     = 0.278 //коэффициент для перевода км/ч в м/с
)

// Walking структура, описывающая тренировку ходьба
type Walking struct {
	Training Training
	Height   float64 //рост пользователя
}

// Calories возвращает количество потраченных килокалорий при ходьбе
// Формула расчета:
// ((0.035 * вес_спортсмена_в_кг + (средняя_скорость_в_метрах_в_секунду**2/рост_в_метрах)
// *0.029 * Вес_спортсмена_в_кг)*время_тренировки_в_часах * мин_в_ч)
// Это переопределенный метод из Calories() из Training
func (w Walking) Calories() float64 {
	cal := CaloriesMeanSpeedMultiplier*w.Training.Weight + math.Pow(w.Training.meanSpeed(), 2)/(w.Height*CmInM)*CaloriesSpeedHeightMultiplier*w.Training.Weight*
		w.Training.Duration.Hours()*MinInHours
	return cal
}

// TrainingInfo возвращает структуру InfoMessage c информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training
func (w Walking) TrainingInfo() InfoMessage {
	im := InfoMessage{
		TrainingType: w.Training.TrainingType,
		Duration:     w.Training.Duration,
		Distance:     w.Training.distance(),
		Speed:        w.Training.meanSpeed(),
		Calories:     w.Training.Calories(),
	}
	return im
}

// Константы для расчета килокалорий, потраченных при плавании
const (
	SwimmingLenStep                  = 1.38 //длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  //коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    //множитель веса пользователя
)

// Swimmming структура, описывающая тренировку Плавание
type Swimming struct {
	Training   Training
	LengthPool int //длина бассейна в метрах
	CountPool  int //количество пересечений бассейна
}

// meanSpeed возвращет среднюю скорость при плавании.
// Формула расчета:
// длина_бассейна * количество_пересечений / м_в_км/продолжительность_тренировки
// Это переопределенный метод Calories() из Training
func (s Swimming) meanSpeed() float64 {
	if math.Abs(s.Training.Duration.Hours()) < eps {
		return 0
	}
	mSpeed := s.LengthPool * s.CountPool / MInKm / int(s.Training.Duration.Hours())
	return float64(mSpeed)
}

// Calories возвращает количество калорий, потраченных при плавании
// Формула расчета (средняя_скорость_в_км/ч + SwimmingCaloriesMeanSpeedShift)*SwimmingCaloriesWeightMultiplier*вес_спортсмена_в_кг * время_тренировки_в_часах
// Это переопределенный метод Calories() из Training()
func (s Swimming) Calories() float64 {
	cal := (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Training.Weight * s.Training.Duration.Hours()
	return cal
}

// TrainingInfo returns info about swimmingTraining
// Это переопределенный метод TrainingInfo() из Training
func (s Swimming) TrainingInfo() InfoMessage {
	im := InfoMessage{
		TrainingType: s.Training.TrainingType,
		Duration:     s.Training.Duration,
		Distance:     s.Training.distance(),
		Speed:        s.Training.meanSpeed(),
		Calories:     s.Training.Calories(),
	}
	return im
}

// ReadData возвращает информацию о проведенной тренировке
func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()
	info.Calories = calories
	return fmt.Sprint(info)
}

func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))
	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}
	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}
	fmt.Println(ReadData(running))
}
