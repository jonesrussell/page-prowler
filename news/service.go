package news

import "fmt"

// Article represents a news article
type Article struct {
	Title       string
	Link        string
	Image       string
	Description string
}

// Service interface defines methods for retrieving news data
type Service interface {
	GetTopStory(siteName string) (Article, error)
	GetBreakingNews(siteName string) ([]Article, error)
	GetLatestUpdates(siteName string) ([]Article, error)
	GetFeatured(siteName string) ([]Article, error)
	GetInPhotos(siteName string) ([]Article, error) // Add this new method
}

// MockService is a simple implementation of Service
type MockService struct{}

func (s *MockService) GetTopStory(siteName string) (Article, error) {
	// Implementation for getting top story
	switch siteName {
	case "cp24":
		return Article{
			Title:       "Ontario Premier Doug Ford says he wants to build a tunnel under Hwy. 401",
			Link:        "https://www.cp24.com/news/ontario-premier-doug-ford-says-he-wants-to-build-a-tunnel-under-hwy-401-1.7051216",
			Image:       "https://www.cp24.com/polopoly_fs/1.7051238.1727270667!/httpImage/image.jpg_gen/derivatives/landscape_620/image.jpg",
			Description: "Premier Doug Ford says he wants to build a tunnel under Highway 401 that would stretch from Brampton to Scarborough.",
		}, nil
	default:
		return Article{}, fmt.Errorf("unknown site: %s", siteName)
	}
}

func (s *MockService) GetBreakingNews(siteName string) ([]Article, error) {
	switch siteName {
	case "cp24":
		return []Article{
			{
				Title:       "Tearful complainant alleges Jacob Hoggard raped, choked her after Hedley concert",
				Link:        "https://www.cp24.com/news/tearful-complainant-alleges-jacob-hoggard-raped-choked-her-after-hedley-concert-1.7051025",
				Image:       "https://www.cp24.com/polopoly_fs/1.7050029.1727194241!/image/image.jpeg_gen/derivatives/landscape_300/image.jpeg",
				Description: "Opening arguments are expected to get underway today in the sexual assault trial of Canadian musician Jacob Hoggard.",
			},
			{
				Title:       "Toronto teachers' union accuses Ford of diverting attention away from Grassy Narrows",
				Link:        "https://www.cp24.com/news/toronto-teachers-union-accuses-ford-of-diverting-attention-away-from-grassy-narrows-as-province-begins-investigating-controversial-field-trip-1.7051645",
				Image:       "https://www.cp24.com/polopoly_fs/1.7051296.1727272549!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Fallout over controversial field trip in Toronto.",
			},
			{
				Title:       "Thieves stole more than $2.2 million of merchandise from moving tractor trailers",
				Link:        "https://www.cp24.com/news/thieves-stole-more-than-2-2-million-of-merchandise-from-moving-tractor-trailers-police-1.7051521",
				Image:       "https://www.cp24.com/polopoly_fs/1.5752494.1648298841!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A Peel Regional Police cruiser is seen in this undated image.",
			},
			{
				Title:       "'This a bright red warning light': Toronto's housing crisis to get worse",
				Link:        "https://www.cp24.com/news/this-a-bright-red-warning-light-toronto-s-housing-crisis-to-get-worse-as-development-applications-drop-off-bild-says-1.7051629",
				Image:       "https://www.cp24.com/polopoly_fs/1.7033313.1726047954!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A new condo construction site is reflected in the window on an ongoing condo construction site in downtown Toronto.",
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown site: %s", siteName)
	}
}

func (s *MockService) GetLatestUpdates(siteName string) ([]Article, error) {
	switch siteName {
	case "cp24":
		return []Article{
			{
				Title:       "'Doug Ford failed the test:' Teachers' union accuses Ford of diverting attention away from Grassy Narrows",
				Link:        "https://www.cp24.com/news/doug-ford-failed-the-test-teachers-union-accuses-ford-of-diverting-attention-away-from-grassy-narrows-1.6851234",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851235.1727280001!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "The Ontario Secondary School Teachers' Federation accuses Premier Doug Ford of diverting attention from the Grassy Narrows mercury poisoning issue.",
			},
			{
				Title:       "'This a bright red warning light': Toronto's housing crisis to get worse as development applications drop off, BILD says",
				Link:        "https://www.cp24.com/news/this-a-bright-red-warning-light-toronto-s-housing-crisis-to-get-worse-as-development-applications-drop-off-bild-says-1.6851236",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851237.1727280002!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "The Building Industry and Land Development Association warns that Toronto's housing crisis could worsen due to a decrease in development applications.",
			},
			{
				Title:       "Woman struck, critically injured by vehicle in Mississauga",
				Link:        "https://www.cp24.com/news/woman-struck-critically-injured-by-vehicle-in-mississauga-1.6851238",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851239.1727280003!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A woman is in critical condition after being struck by a vehicle in Mississauga.",
			},
			{
				Title:       "Bloc gives Liberals Oct. 29 deadline to meet demands or face potential early election",
				Link:        "https://www.cp24.com/news/bloc-gives-liberals-oct-29-deadline-to-meet-demands-or-face-potential-early-election-1.6851240",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851241.1727280004!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "The Bloc Québécois sets an October 29 deadline for the Liberal government to meet their demands or potentially face an early election.",
			},
			{
				Title:       "Tropical Storm Helene brings wet weather to Toronto",
				Link:        "https://www.cp24.com/news/tropical-storm-helene-brings-wet-weather-to-toronto-1.6851242",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851243.1727280005!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Tropical Storm Helene brings rainy weather to Toronto and surrounding areas.",
			},
			{
				Title:       "Hamilton police searching for additional suspects in Grindr ambush robberies",
				Link:        "https://www.cp24.com/news/hamilton-police-searching-for-additional-suspects-in-grindr-ambush-robberies-1.6851244",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851245.1727280006!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Hamilton police are looking for more suspects involved in a series of ambush robberies targeting Grindr users.",
			},
			{
				Title:       "Statistics Canada says population grew 0.6 per cent in Q2 to 41,288,599",
				Link:        "https://www.cp24.com/news/statistics-canada-says-population-grew-0-6-per-cent-in-q2-to-41-288-599-1.6851246",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851247.1727280007!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Statistics Canada reports a 0.6 percent population growth in the second quarter, bringing the total to 41,288,599.",
			},
			{
				Title:       "Drought in Brazil, Vietnam highlight climate change's impact on coffee: experts",
				Link:        "https://www.cp24.com/news/drought-in-brazil-vietnam-highlight-climate-change-s-impact-on-coffee-experts-1.6851248",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851249.1727280008!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Experts warn that droughts in Brazil and Vietnam are highlighting the impact of climate change on coffee production.",
			},
			{
				Title:       "Ministry of Education launches probe into TDSB field trip to rally",
				Link:        "https://www.cp24.com/news/ministry-of-education-launches-probe-into-tdsb-field-trip-to-rally-1.6851250",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851251.1727280009!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "The Ontario Ministry of Education initiates an investigation into a Toronto District School Board field trip to a rally.",
			},
			{
				Title:       "Toronto jazz musician fatally struck in collision remembered as 'talented,' 'beautiful' person",
				Link:        "https://www.cp24.com/news/toronto-jazz-musician-fatally-struck-in-collision-remembered-as-talented-beautiful-person-1.6851252",
				Image:       "https://www.cp24.com/polopoly_fs/1.6851253.1727280010!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A Toronto jazz musician who died in a collision is remembered by the community as a talented and beautiful person.",
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown site: %s", siteName)
	}
}

func (s *MockService) GetFeatured(siteName string) ([]Article, error) {
	if siteName != "cp24" {
		return nil, fmt.Errorf("unknown site: %s", siteName)
	}

	return []Article{
		{
			Title:       "'I'm here for the Porsche': Video shows brazen car theft in Mississauga",
			Link:        "https://www.cp24.com/news/i-m-here-for-the-porsche-video-shows-brazen-car-theft-in-mississauga-1.7042881",
			Image:       "https://www.cp24.com/polopoly_fs/1.7042887.1726680149!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "Police are searching for a suspect who allegedly stole a luxury SUV in Mississauga that was listed for sale on Auto Trader.",
		},
		{
			Title:       "Samfiru Tumarkin takes the stress out of legal fees with no upfront costs",
			Link:        "https://www.cp24.com/news/samfiru-tumarkin-takes-the-stress-out-of-legal-fees-with-no-upfront-costs-1.7050093",
			Image:       "https://www.cp24.com/polopoly_fs/1.7050133.1727198597!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "SPONSORED: Samfiru Tumarkin offers legal services with no upfront costs.",
		},
		{
			Title:       "Video shows suspects running down street after allegedly setting St. Catharines restaurant on fire",
			Link:        "https://www.cp24.com/news/video-shows-suspects-running-down-street-after-allegedly-setting-st-catharines-restaurant-on-fire-1.7040843",
			Image:       "https://www.cp24.com/polopoly_fs/1.7040849.1726573859!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "Two suspects are wanted in connection with an arson investigation in St. Catharines.",
		},
	}, nil
}

func (s *MockService) GetInPhotos(siteName string) ([]Article, error) {
	if siteName != "cp24" {
		return nil, fmt.Errorf("unknown site: %s", siteName)
	}

	return []Article{
		{
			Title:       "Stars descend on Toronto for TIFF 2024",
			Link:        "https://www.cp24.com/photo-galleries/stars-descend-on-toronto-for-tiff-2024-1.7027599",
			Image:       "https://www.cp24.com/polopoly_fs/1.7038892.1726435536!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "Celebrities are appearing on the streets of Toronto as TIFF gets underway.",
		},
		{
			Title:       "Fan Expo Canada 2024",
			Link:        "https://www.cp24.com/fan-expo-canada-2024-1.7012679",
			Image:       "https://www.cp24.com/polopoly_fs/1.7015460.1724759164!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "Fans, creators and celebrities come together at the massive annual event in Toronto.",
		},
		{
			Title:       "High-profile cases that gripped Toronto",
			Link:        "https://www.cp24.com/high-profile-cases-that-gripped-toronto-1.6969783",
			Image:       "https://www.cp24.com/polopoly_fs/1.6969831.1721357494!/httpImage/image.jpg_gen/derivatives/landscape_800/image.jpg",
			Description: "A look at some of the high-profile cases in the Greater Toronto Area that gripped the city in recent years.",
		},
	}, nil
}

// NewMockService creates a new instance of MockService
func NewMockService() Service {
	return &MockService{}
}
