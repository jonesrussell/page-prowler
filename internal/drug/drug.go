package drug

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

func Related(href string) bool {
	// We only want the title part of the url
	sliced := strings.Split(href, "/")
	title := sliced[len(sliced)-1]

	// Remove -'s from title
	sliced = strings.Split(title, "-")
	title = strings.Join(sliced, " ")

	// Remove stopwords
	title = stopwords.CleanString(title, "en", false)

	// Trim
	title = strings.TrimSpace(title)

	// Stem the remaining words
	sliced = strings.Split(title, " ")
	sliced = stemmer.StemMultiple(sliced)

	// Lemmatize
	/*lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(slicedHref); i++ {
		slicedHref[i] = lemmatizer.Lemma(slicedHref[i])
	}*/

	// Convert slice back to string
	title = strings.Join(sliced, " ")

	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	similarityDrug := strutil.Similarity("DRUG", title, swg)
	similaritySmokeJoint := strutil.Similarity("SMOKE JOINT", title, swg)
	similarityGrowop := strutil.Similarity("GROW OP", title, swg)
	similarityCannabi := strutil.Similarity("CANNABI", title, swg)
	similarityImpair := strutil.Similarity("IMPAIR", title, swg)
	similarityShoot := strutil.Similarity("SHOOT", title, swg)
	similarityFirearm := strutil.Similarity("FIREARM", title, swg)
	similarityMurder := strutil.Similarity("MURDER", title, swg)
	similarityCocain := strutil.Similarity("COCAIN", title, swg)

	return similarityDrug == 1 ||
		similaritySmokeJoint == 1 ||
		similarityGrowop == 1 ||
		similarityCannabi == 1 ||
		similarityImpair == 1 ||
		similarityShoot == 1 ||
		similarityFirearm == 1 ||
		similarityMurder == 1 ||
		similarityCocain == 1
}
