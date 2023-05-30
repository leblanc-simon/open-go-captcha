package captcha

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hashicorp/go-uuid"
	"golang.org/x/exp/slices"
	"leblanc.io/open-go-captcha/config"
	"leblanc.io/open-go-captcha/crypto"
)

var numberOfAvailableIcons int = 100

type Icon struct {
	DataURI string
	Identifier string
}

type Captcha struct {
	Session string
	Token string 
	Icons []Icon
	Answers []string
}

// @see calculateIconAmounts from https://github.com/fabianwennink/IconCaptcha-Plugin-jQuery-PHP/blob/master/src/captcha.class.php#L611
func calculateIconAmounts(nbIcons int, correctNbIcons int) []int {
	remainder := nbIcons - correctNbIcons
	remainderDivided := float64(remainder) / float64(2)
	pickDivided := (rand.Intn(4) + 1) != 1

	if math.Mod(remainderDivided, 1) != 0.0 && pickDivided {
		left := math.Floor(remainderDivided)
		right := math.Ceil(remainderDivided)

		if left > float64(correctNbIcons) && right > float64(correctNbIcons) {
			return []int{int(left), int(right)}
		}
	} else if pickDivided && int(remainderDivided) > correctNbIcons {
		return []int{int(remainderDivided), int(remainderDivided)}
	}

	return []int{int(remainder)}
}

func getIconData(c *config.Config) ([]Icon, []string, error) {
	rand.Seed(time.Now().UnixNano())

	nbIcons := rand.Intn(c.Captcha.Max - c.Captcha.Min + 1) + c.Captcha.Min
	correctNbIcons := rand.Intn(c.Captcha.MaxGood) + 1

	iconsAmount := calculateIconAmounts(nbIcons, correctNbIcons)
	iconsAmount = append(iconsAmount, correctNbIcons)

	var icons []Icon
	var answers []string

	for _, iconAmount := range iconsAmount {
		icon := getRandomIcon(c, icons)
		
		for iterator := 1; iterator <= iconAmount; iterator++ {
			identifier, err := uuid.GenerateUUID()
			if err != nil {
				fmt.Println(err)
				return nil, nil, err
			}

			icons = append(icons, Icon{
				DataURI: icon,
				Identifier: identifier,
			})

			if iconAmount == correctNbIcons {
				answers = append(answers, identifier)
			}
		}
	}

	rand.Shuffle(len(icons), func(i int, j int) {
		icons[i], icons[j] = icons[j], icons[i]
	})

	return icons, answers, nil
}

func getRandomIcon(c *config.Config, currentIcons []Icon) (string) {
	find := false

	for {
		icon := fmt.Sprintf("%d.svg", rand.Intn(numberOfAvailableIcons))

		if !slices.ContainsFunc(currentIcons, func (searchIcon Icon) bool { return searchIcon.DataURI == icon }) {
			find = true
		}

		if find {
			return icon
		}
	}
}

func NewCaptcha(c *config.Config) (*Captcha, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	token := crypto.Encrypt(id)

	iconData, anwsers, err := getIconData(c)
	if err != nil {
		return nil, err
	}
	
	fmt.Println(id, iconData, anwsers)

	return &Captcha{
		Session: id,
		Token: token,
		Icons: iconData,
		Answers: anwsers,
	}, nil
}