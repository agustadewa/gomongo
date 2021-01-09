package tools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agustadewa/gomongo"
	"image/png"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/jung-kurt/gofpdf"
	validatorv9 "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// Tools type
type Tools struct{}

// PrintPDF method
func (tool Tools) PrintPDF(name, callSign, band, templatePath, outPath, fileType string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("ArchivoBlack-Regular", "", "ArchivoBlack-Regular.json")
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("ATOMICCLOCKRADIO", "", "ATOMICCLOCKRADIO.json")

	pdf.SetHeaderFunc(func() {
		// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
		pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

		pdf.SetFont("ArchivoBlack-Regular", "", 47)
		pdf.SetXY(4, 91)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(40, 10, callSign)

		pdf.SetFont("ArchivoBlack-Regular", "", 25)
		pdf.SetXY(6, 105)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(10, 10, name)

		pdf.SetFont("ATOMICCLOCKRADIO", "", 23)
		pdf.SetTextColor(255, 255, 255)
		if band == "40 m" {
			pdf.SetXY(131, 43)
			pdf.Cell(10, 10, "7.135")
		} else if band == "2 m" {
			pdf.SetXY(119, 43)
			pdf.Cell(10, 10, "145.240")
		}
	})

	err := pdf.OutputFileAndClose(outPath)
	if err != nil {
		fmt.Println("ERRRORRR! ", err)
	}
	return err
}

// PrintPDFV2 method
func (tool Tools) PrintPDFV2(name, callSign, band, templatePath, fileType string, w io.Writer) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("OrangeTypewriter", "", "OrangeTypewriter.json")
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("Kanit-Bold", "", "Kanit-Bold.json")

	pdf.SetHeaderFunc(func() {
		// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
		pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

		// CALL SIGN
		pdf.SetFont("Kanit-Bold", "", 48)
		pdf.SetXY(247, 80)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(40, 10, callSign, "", 0, "R", false, 0, "")

		// NAME
		pdf.SetFont("Kanit-Bold", "", 18)
		pdf.SetXY(276, 95)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(10, 10, name, "", 0, "R", false, 0, "")

		// FREQUENCY
		pdf.SetFont("OrangeTypewriter", "", 16)
		pdf.SetTextColor(0, 0, 0)
		if band == "40 m" {
			pdf.SetXY(279, 23)
			pdf.CellFormat(10, 10, "7.135 MHz", "", 0, "R", false, 0, "")
		} else if band == "2 m" {
			pdf.SetXY(279, 23)
			pdf.CellFormat(10, 10, "145.240 MHz", "", 0, "R", false, 0, "")
		}
	})

	err := pdf.Output(w)
	if err != nil {
		log.Println("error creating pdf:", err)
	}
	return err
}

// PrintPDFV3 method
func (tool Tools) PrintPDFV3(name, callSign, band, frequency, templatePath, fileType string, w io.Writer, imageCertTemplate gomongo.ImageCertTemplate) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFontLocation(imageCertTemplate.TemplateProperties.CallSign.FontDir)
	pdf.AddFont(imageCertTemplate.TemplateProperties.CallSign.FontName, "", fmt.Sprintf("%s.json", imageCertTemplate.TemplateProperties.CallSign.FontName))
	pdf.SetFontLocation(imageCertTemplate.TemplateProperties.IdentityName.FontDir)
	pdf.AddFont(imageCertTemplate.TemplateProperties.IdentityName.FontName, "", fmt.Sprintf("%s.json", imageCertTemplate.TemplateProperties.IdentityName.FontName))
	pdf.SetFontLocation(imageCertTemplate.TemplateProperties.Frequency.FontDir)
	pdf.AddFont(imageCertTemplate.TemplateProperties.Frequency.FontName, "", fmt.Sprintf("%s.json", imageCertTemplate.TemplateProperties.Frequency.FontName))
	handler := func(imageCertTemplate gomongo.ImageCertTemplate) func() {
		return func() {
			// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
			pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

			// CALL SIGN
			pdf.SetFont(imageCertTemplate.TemplateProperties.CallSign.FontName, "", imageCertTemplate.TemplateProperties.CallSign.FontSize)
			pdf.SetXY(imageCertTemplate.TemplateProperties.CallSign.TextPosition.X, imageCertTemplate.TemplateProperties.CallSign.TextPosition.Y)
			pdf.SetTextColor(imageCertTemplate.TemplateProperties.CallSign.FontColor.R, imageCertTemplate.TemplateProperties.CallSign.FontColor.G, imageCertTemplate.TemplateProperties.CallSign.FontColor.B)
			pdf.CellFormat(40, 10, callSign, "", 0, imageCertTemplate.TemplateProperties.CallSign.TextAlign, false, 0, "")

			// NAME
			pdf.SetFont(imageCertTemplate.TemplateProperties.IdentityName.FontName, "", imageCertTemplate.TemplateProperties.IdentityName.FontSize)
			pdf.SetXY(imageCertTemplate.TemplateProperties.IdentityName.TextPosition.X, imageCertTemplate.TemplateProperties.IdentityName.TextPosition.Y)
			pdf.SetTextColor(imageCertTemplate.TemplateProperties.IdentityName.FontColor.R, imageCertTemplate.TemplateProperties.IdentityName.FontColor.G, imageCertTemplate.TemplateProperties.IdentityName.FontColor.B)
			pdf.CellFormat(10, 10, name, "", 0, imageCertTemplate.TemplateProperties.IdentityName.TextAlign, false, 0, "")

			// FREQUENCY
			pdf.SetFont(imageCertTemplate.TemplateProperties.Frequency.FontName, "", imageCertTemplate.TemplateProperties.Frequency.FontSize)
			pdf.SetTextColor(imageCertTemplate.TemplateProperties.Frequency.FontColor.R, imageCertTemplate.TemplateProperties.Frequency.FontColor.G, imageCertTemplate.TemplateProperties.Frequency.FontColor.B)
			pdf.SetXY(imageCertTemplate.TemplateProperties.Frequency.TextPosition.X, imageCertTemplate.TemplateProperties.Frequency.TextPosition.Y)
			pdf.CellFormat(10, 10, fmt.Sprintf("%s - %s", frequency, band), "", 0, imageCertTemplate.TemplateProperties.Frequency.TextAlign, false, 0, "")
		}
	}
	pdf.SetHeaderFunc(handler(imageCertTemplate))

	err := pdf.Output(w)
	if err != nil {
		log.Println("error creating pdf:", err)
	}
	return err
}

// SaveImageFromB64 method
func (tool Tools) SaveImageFromB64(b64 string, filePath string) error {
	unbased, errDecode := base64.StdEncoding.Strict().DecodeString(b64)
	if errDecode != nil {
		fmt.Println(errDecode)
	}

	reader := bytes.NewReader(unbased)

	// Decode JPG
	//img, errDecodeJpeg := jpeg.Decode(reader);
	//if  errDecodeJpeg != nil {
	//	panic("BAD JPG")
	//}

	// Decode PNG
	img, errDecodePng := png.Decode(reader)
	if errDecodePng != nil {
		panic("BAD PNG")
	}

	// Create File
	file, errCreateFile := os.Create(filePath)
	if errCreateFile != nil {
		panic(errCreateFile)
	}
	defer file.Close()

	// Encode JPG
	//jpgOpt := jpeg.Options{Quality: 30}
	//errEncodeJpg = jpeg.Encode(f, img, &jpgOpt)

	// Encode PNG
	errEncodePng := png.Encode(file, img)
	if errEncodePng != nil {
		fmt.Println(errEncodePng)
	}

	//return errEncodeJpg
	return errEncodePng
}

// ReplaceRegex method
func (tool Tools) ReplaceRegex(b64url *string, b64 *string) {
	var re = regexp.MustCompile(`^[^,]*,`)
	*b64 = re.ReplaceAllString(*b64url, "")
}

// RNDString method
func (tool Tools) RNDString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// BindValidate function
func BindValidate(o interface{}) error {
	ginValidate, ok := binding.Validator.Engine().(*validatorv9.Validate)
	if !ok {
		err := binding.Validator.ValidateStruct(o)
		if err != nil {
			return err
		}
		validate = validator.New()
		err = validate.Struct(o)

		return err
	}

	err := ginValidate.Struct(o)
	if err != nil {
		return err
	}

	validate = validator.New()
	err = validate.Struct(o)

	return err
}

// PairValues function
func PairValues(i, o interface{}) error {
	if i == nil {
		return errors.New("error while pair values, values is nil")
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &o)
	if err != nil {
		return err
	}

	//check type of o
	r := reflect.ValueOf(o)
	if r.Kind() == reflect.Ptr && !r.IsNil() {
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct && r.Kind() != reflect.Interface {

		return nil
	}

	//validate struct :
	err = BindValidate(o)
	if err != nil {

		return err
	}

	return nil
}
