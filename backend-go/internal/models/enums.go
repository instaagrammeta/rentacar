package models

// Enumerations used across the domain models. The string value is the canonical
// machine value stored in the database (VARCHAR). Human readable Russian labels
// are exposed through the Label methods and used for exports (Excel / PDF).

// UserRole controls access to modules.
type UserRole string

const (
	RoleAdministrator UserRole = "administrator"
	RoleRentalManager UserRole = "rental_manager"
	RoleCashier       UserRole = "cashier"
	RoleOperator      UserRole = "operator"
)

var userRoleLabels = map[UserRole]string{
	RoleAdministrator: "Администратор",
	RoleRentalManager: "Менеджер аренды",
	RoleCashier:       "Кассир",
	RoleOperator:      "Оператор",
}

// Label returns the Russian display label for the role.
func (r UserRole) Label() string { return userRoleLabels[r] }

// Valid reports whether the role is a known value.
func (r UserRole) Valid() bool { _, ok := userRoleLabels[r]; return ok }

// AllUserRoles returns every role in declaration order.
func AllUserRoles() []UserRole {
	return []UserRole{RoleAdministrator, RoleRentalManager, RoleCashier, RoleOperator}
}

// ClientStatus represents the verification state of a client.
type ClientStatus string

const (
	ClientActive              ClientStatus = "active"
	ClientPendingVerification ClientStatus = "pending_verification"
	ClientBlacklisted         ClientStatus = "blacklisted"
)

var clientStatusLabels = map[ClientStatus]string{
	ClientActive:              "Активен",
	ClientPendingVerification: "Ожидает проверки",
	ClientBlacklisted:         "В чёрном списке",
}

func (s ClientStatus) Label() string { return clientStatusLabels[s] }
func (s ClientStatus) Valid() bool   { _, ok := clientStatusLabels[s]; return ok }

// CarStatus represents the availability of a vehicle.
type CarStatus string

const (
	CarAvailable   CarStatus = "available"
	CarReserved    CarStatus = "reserved"
	CarRented      CarStatus = "rented"
	CarMaintenance CarStatus = "maintenance"
)

var carStatusLabels = map[CarStatus]string{
	CarAvailable:   "Доступен",
	CarReserved:    "Забронирован",
	CarRented:      "В аренде",
	CarMaintenance: "На обслуживании",
}

func (s CarStatus) Label() string { return carStatusLabels[s] }
func (s CarStatus) Valid() bool   { _, ok := carStatusLabels[s]; return ok }

func AllCarStatuses() []CarStatus {
	return []CarStatus{CarAvailable, CarReserved, CarRented, CarMaintenance}
}

// ReservationStatus represents the state of a reservation.
type ReservationStatus string

const (
	ReservationReserved  ReservationStatus = "reserved"
	ReservationConfirmed ReservationStatus = "confirmed"
	ReservationCancelled ReservationStatus = "cancelled"
)

var reservationStatusLabels = map[ReservationStatus]string{
	ReservationReserved:  "Забронировано",
	ReservationConfirmed: "Подтверждено",
	ReservationCancelled: "Отменено",
}

func (s ReservationStatus) Label() string { return reservationStatusLabels[s] }
func (s ReservationStatus) Valid() bool   { _, ok := reservationStatusLabels[s]; return ok }

// RentalStatus represents the state of a rental contract.
type RentalStatus string

const (
	RentalActive    RentalStatus = "active"
	RentalCompleted RentalStatus = "completed"
	RentalCancelled RentalStatus = "cancelled"
)

var rentalStatusLabels = map[RentalStatus]string{
	RentalActive:    "Активна",
	RentalCompleted: "Завершена",
	RentalCancelled: "Отменена",
}

func (s RentalStatus) Label() string { return rentalStatusLabels[s] }
func (s RentalStatus) Valid() bool   { _, ok := rentalStatusLabels[s]; return ok }

func AllRentalStatuses() []RentalStatus {
	return []RentalStatus{RentalActive, RentalCompleted, RentalCancelled}
}

// PaymentMethod represents how a payment was made.
type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "cash"
	PaymentBankTransfer PaymentMethod = "bank_transfer"
	PaymentCard         PaymentMethod = "card"
)

var paymentMethodLabels = map[PaymentMethod]string{
	PaymentCash:         "Наличные",
	PaymentBankTransfer: "Банковский перевод",
	PaymentCard:         "Карта",
}

func (m PaymentMethod) Label() string { return paymentMethodLabels[m] }
func (m PaymentMethod) Valid() bool   { _, ok := paymentMethodLabels[m]; return ok }

// PaymentType represents the purpose of a payment.
type PaymentType string

const (
	PaymentTypeDeposit PaymentType = "deposit"
	PaymentTypeRental  PaymentType = "rental"
	PaymentTypePenalty PaymentType = "penalty"
	PaymentTypeDamage  PaymentType = "damage"
	PaymentTypeRefund  PaymentType = "refund"
)

var paymentTypeLabels = map[PaymentType]string{
	PaymentTypeDeposit: "Депозит",
	PaymentTypeRental:  "Оплата аренды",
	PaymentTypePenalty: "Штраф",
	PaymentTypeDamage:  "Возмещение ущерба",
	PaymentTypeRefund:  "Возврат средств",
}

func (t PaymentType) Label() string { return paymentTypeLabels[t] }
func (t PaymentType) Valid() bool   { _, ok := paymentTypeLabels[t]; return ok }

// BlacklistReason represents why a client was blacklisted.
type BlacklistReason string

const (
	ReasonFraud            BlacklistReason = "fraud"
	ReasonVehicleDamage    BlacklistReason = "vehicle_damage"
	ReasonNonPayment       BlacklistReason = "non_payment"
	ReasonSeriousViolation BlacklistReason = "serious_violation"
)

var blacklistReasonLabels = map[BlacklistReason]string{
	ReasonFraud:            "Мошенничество",
	ReasonVehicleDamage:    "Повреждение автомобиля",
	ReasonNonPayment:       "Неоплата",
	ReasonSeriousViolation: "Серьёзное нарушение",
}

func (r BlacklistReason) Label() string { return blacklistReasonLabels[r] }
func (r BlacklistReason) Valid() bool   { _, ok := blacklistReasonLabels[r]; return ok }

func AllBlacklistReasons() []BlacklistReason {
	return []BlacklistReason{ReasonFraud, ReasonVehicleDamage, ReasonNonPayment, ReasonSeriousViolation}
}
