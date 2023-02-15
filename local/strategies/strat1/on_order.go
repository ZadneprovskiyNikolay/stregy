package strat1

import (
	"stregy/internal/domain/order"
	"stregy/local/strategies/posmerge"
)

/*
scenarios:

	filled:
		new main order - merge
		other:
			sl order - nothing(delete current mergedPosition if closed in OnMergedPositionOrder)
			tp order - nothing(delete current mergedPosition if closed in OnMergedPositionOrder)
	rejected:
	    new main order - report
		other:
		    sl order - nothing(report in OnMergedPositionOrder)
			tp order - nothing(report in OnMergedPositionOrder)
*/
func (s *strat1) OnOrder(o order.Order) {
	// s.broker.Printf("OnOrder: %v", o)

	s.positionMerger.OnOrder(o)

	if s.newMainOrder != nil {

		if o.ID == s.newMainOrder.ID {

			if o.Status == order.FilledOrder {
				var mergedPosition posmerge.MergedPosition
				var err error
				if s.mergedPosition == nil {
					mergedPosition, err = s.positionMerger.MergePositions(o.Position)
					if err != nil {
						panic(err)
					}
					s.totalAverages = 0
					s.nextAveragingPrice = getNextAveragingPrice(o.Diraction, o.ExecutionPrice)
				} else {
					mergedPosition, err = s.positionMerger.MergePosition(o.Position, s.mergedPosition.ID)
					if err != nil {
						panic(err)
					}
					s.positionMerger.SetTakeProfit(mergedPosition.ID, s.newTpPrice)
					s.positionMerger.SetStopLoss(mergedPosition.ID, s.newSlPrice)
					s.totalAverages += 1
					s.nextAveragingPrice = getNextAveragingPrice(o.Diraction, s.nextAveragingPrice)
				}
				if err != nil {
					s.broker.Printf("failed to merge positions: %v", err)
					panic(err)
				}

				s.mergedPosition = &mergedPosition
				s.newMainOrder = nil
				s.newSlPrice = 0
				s.newTpPrice = 0

			} else if o.Status == order.RejectedOrder {
				s.broker.Printf("main order %v rejected", o)
			}

		} else if o.Status == order.RejectedOrder {
			s.broker.Printf("ctg order rejected: %v", o)
		}
	}
}
