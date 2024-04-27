package payme

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Builder struct {
	iban             string
	amount           int
	paymentReference string
	specificSymbol   string
	note             string
}

const (
	versionKey               = "V"
	ibanKey                  = "IBAN"
	amountKey                = "AM"
	currencyCodeKey          = "CC"
	paymentIdentificationKey = "PI"
	messageKey               = "MSG"
	creditorsKey             = "CN"
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) IBAN(val string) *Builder {
	b.iban = val
	return b
}

func (b *Builder) Amount(val int) *Builder {
	b.amount = val
	return b
}

func (b *Builder) PaymentReference(val string) *Builder {
	b.paymentReference = val
	return b
}

func (b *Builder) SpecificSymbol(val string) *Builder {
	b.specificSymbol = val
	return b
}

func (b *Builder) Note(val string) *Builder {
	b.note = val
	return b
}

func (b *Builder) Build() (string, error) {
	v := url.Values{}
	v.Set(versionKey, "1")
	v.Set(ibanKey, strings.ReplaceAll(b.iban, " ", ""))
	v.Set(amountKey, strconv.Itoa(b.amount))
	v.Set(currencyCodeKey, "EUR")
	v.Set(paymentIdentificationKey, fmt.Sprintf("/VS%s/SS%s/KS%s", b.paymentReference, b.specificSymbol, ""))
	v.Set(creditorsKey, "salezko")
	if b.note != "" {
		v.Set(messageKey, regexp.MustCompile(`[^\p{L}\p{N} ]+`).ReplaceAllString(b.note, ""))
	}
	return "https://payme.sk?" + v.Encode(), nil

}
